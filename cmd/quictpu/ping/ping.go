package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/text"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"path"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"github.com/mr-tron/base58"
	"k8s.io/klog/v2"
)

// QUIC_GO_LOG_LEVEL=DEBUG

var (
	flagDebug      = flag.Bool("debug", false, "Enable debug logging")
	flagCount      = flag.Int("c", 1, "Number of pings to send, -1 for infinite")
	flagDelay      = flag.Duration("i", 1*time.Second, "Delay between pings")
	flagAddr       = flag.String("addr", "", "Address to ping (<host>:<port>)")
	flagSourcePort = flag.Int("s", 0, "Source port to use (0 for random/default)")
	flagKey        = flag.String("k", "", "Path to private key file (default ~/.config/solana/id.json)")
	flagSendTx     = flag.Bool("send-tx", false, "Send a transaction")
	flagRPC        = flag.String("u", "", "RPC URL to use for getting blockhash")
)

type pingData struct {
	Slot  uint64    `json:"slot"`
	Ts    time.Time `json:"ts"`
	Index int       `json:"index"`
}

func init() {
	klog.InitFlags(nil)
	flag.Parse()

	if *flagKey == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		*flagKey = path.Join(home, ".config/solana/id.json")
	}

	if *flagSendTx && *flagRPC == "" {
		klog.Exitf("RPC URL must be specified when sending transactions")
	}

	// Mute receive buffer warning (we don't even send data!)
	if err := os.Setenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING", "1"); err != nil {
		panic(err)
	}
}

func loadLocalSigner() (solana.PrivateKey, error) {
	return solana.PrivateKeyFromSolanaKeygenFile(*flagKey)
}

func generateRandomX509KeyPair() (tls.Certificate, error) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return tls.Certificate{}, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	notBefore := time.Now().Add(-1 * time.Hour)
	notAfter := notBefore.Add(1 * time.Hour)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Radiance Certificate Manufacturing Co."},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  priv,
	}, nil
}

func buildTransaction(now time.Time, i int, blockhash solana.Hash, feePayer solana.PublicKey) *solana.Transaction {
	payload := &pingData{Ts: now, Index: i}
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	ins := solana.NewInstruction(solana.MemoProgramID, solana.AccountMetaSlice{}, b)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{ins}, blockhash, solana.TransactionPayer(feePayer))
	if err != nil {
		panic(err)
	}

	return tx
}

func main() {
	ctx := context.Background()

	if *flagAddr == "" {
		klog.Exit("No address to ping specified")
	}

	var (
		qconf quic.Config
		dbg   io.Writer
		err   error
	)
	if *flagDebug {
		dbg, err = os.OpenFile("keylog.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			klog.Exitf("Failed to open keylog file: %v", err)
		}
		qconf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			filename := fmt.Sprintf("client_%x.qlog", connID)
			f, err := os.Create(filename)
			if err != nil {
				klog.Fatal(err)
			}
			log.Printf("Creating qlog file %s.\n", filename)
			return f
		})
	}

	clientCert, err := generateRandomX509KeyPair()
	if err != nil {
		panic(err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"solana-tpu"},
		KeyLogWriter:       dbg,
		Certificates:       []tls.Certificate{clientCert},
	}

	for c := 0; c < *flagCount || *flagCount == -1; c++ {
		c := c
		t := time.Now()
		minTimeout := 1 * time.Second
		if *flagDelay > minTimeout {
			minTimeout = *flagDelay
		}
		ctx, cancel := context.WithTimeout(ctx, minTimeout)

		udpAddr, err := net.ResolveUDPAddr("udp", *flagAddr)
		if err != nil {
			klog.Exitf("Failed to resolve UDP address: %v", err)
		}
		udpConn, err := net.ListenUDP("udp",
			&net.UDPAddr{IP: net.IPv4zero, Port: *flagSourcePort})
		if err != nil {
			klog.Exitf("Failed to listen on UDP socket: %v", err)
		}

		conn, err := quic.DialContext(ctx, udpConn, udpAddr, *flagAddr, tlsConf, &qconf)
		if err != nil {
			klog.Errorf("Failed to dial: %v", err)
			time.Sleep(*flagDelay)
			cancel()
			continue
		}
		cancel()

		klog.Infof("Connected to %s (in %dms, %d/%d)",
			*flagAddr, time.Since(t).Milliseconds(),
			c+1, *flagCount)

		if klog.V(1).Enabled() {
			for _, cert := range conn.ConnectionState().TLS.PeerCertificates {
				klog.Infof("Certificate: %s", cert.Subject)
				klog.Infof("Public key: %s", base58.Encode(cert.PublicKey.(ed25519.PublicKey)))
			}
		}

		if *flagSendTx {
			client := rpc.New(*flagRPC)
			// TODO: close

			out, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
			if err != nil {
				klog.Exitf("Failed to get recent blockhash: %v", err)
			}

			signer, err := loadLocalSigner()
			if err != nil {
				klog.Exitf("Failed to load local signer: %v", err)
			}

			tx := buildTransaction(t, c, out.Value.Blockhash, signer.PublicKey())
			sigs, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
				if key != signer.PublicKey() {
					panic("no private key for unknown signer " + key.String())
				}
				return &signer
			})
			if err != nil {
				panic(err)
			}

			if klog.V(2).Enabled() {
				tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Ping memo"))
			}

			txb, err := tx.MarshalBinary()
			if err != nil {
				panic(err)
			}

			klog.Infof("Sending tx %s", sigs[0].String())
			klog.V(2).Infof("tx: %s", hex.EncodeToString(txb))

			// Open a stream
			stream, err := conn.OpenUniStream()
			if err != nil {
				klog.Errorf("Failed to open stream: %v", err)
				continue
			}

			if n, err := stream.Write(txb); err != nil {
				klog.Errorf("Failed to write to stream: %v", err)
				continue
			} else {
				klog.V(2).Infof("Wrote %d bytes to stream", n)
			}

			if err := stream.Close(); err != nil {
				klog.Errorf("Failed to close stream: %v", err)
				continue
			}
		}

		if err := conn.CloseWithError(0, ""); err != nil {
			klog.Exitf("Failed to close: %v", err)
		}

		time.Sleep(*flagDelay)
	}
}
