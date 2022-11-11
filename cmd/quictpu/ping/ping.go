package main

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"github.com/mr-tron/base58"
	"k8s.io/klog/v2"
)

// QUIC_GO_LOG_LEVEL=DEBUG

var (
	flagDebug = flag.Bool("debug", false, "Enable debug logging")
	flagCount = flag.Int("c", 1, "Number of pings to send, -1 for infinite")
	flagDelay = flag.Duration("i", 1*time.Second, "Delay between pings")
	flagAddr  = flag.String("addr", "", "Address to ping (<host>:<port>)")
)

func init() {
	klog.InitFlags(nil)
	flag.Parse()

	// Mute receive buffer warning (we don't even send data!)
	if err := os.Setenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING", "1"); err != nil {
		panic(err)
	}
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

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"solana-tpu"},
		KeyLogWriter:       dbg,
	}

	for c := 0; c < *flagCount || *flagCount == -1; c++ {
		c := c
		t := time.Now()
		minTimeout := 100 * time.Millisecond
		if *flagDelay > minTimeout {
			minTimeout = *flagDelay
		}
		ctx, cancel := context.WithTimeout(ctx, minTimeout)
		conn, err := quic.DialAddrContext(ctx, *flagAddr, tlsConf, &qconf)
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

		if err := conn.CloseWithError(0, ""); err != nil {
			klog.Exitf("Failed to close: %v", err)
		}

		time.Sleep(*flagDelay)
	}
}
