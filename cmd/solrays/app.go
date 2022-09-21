package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"k8s.io/klog/v2"

	"go.firedancer.io/radiance/pkg/util"
)

var (
	backendAddr = flag.String("backend", "http://127.0.0.1:8899", "Backend RPC URI to proxy to")
	debugAddr   = flag.String("listen", "[::1]:6060", "pprof and metrics server address")
	tlsHostname = flag.String("tlsHostname", "", "When set, serve TLS using Let's Encrypt using the hostname in question")
	tlsProd     = flag.Bool("tlsProd", false, "Use the production Let's Encrypt environment")
	cacheDir    = flag.String("cacheDir", "solrays-data", "Cache directory")
)

func init() {
	klog.CopyStandardLogTo("INFO")
	klog.InitFlags(nil)
	flag.Parse()
}

func getSDListeners() []net.Listener {
	// We use systemd socket activation for (almost) zero downtime deployment -
	// systemd will keep the socket open even while we restart the process
	// (plus, it allows us to bind to port 80).
	//
	// Read more: https://vincent.bernat.ch/en/blog/2018-systemd-golang-socket-activation

	listeners, err := activation.Listeners()
	if err != nil {
		klog.Fatalf("cannot retrieve listeners: %s", err)
	}
	if len(listeners) != 1 {
		klog.Fatalf("unexpected number of sockets passed by systemd (%d != 1)", len(listeners))
	}

	return listeners
}

func shutdownHandler(server *http.Server) chan struct{} {
	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		klog.Info("server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			klog.Exitf("cannot gracefully shut down the server: %s", err)
		}
		close(done)
	}()
	return done
}

func main() {
	listeners := getSDListeners()

	// Metrics recording middleware
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	mux := http.NewServeMux()

	mux.Handle("/", newHandler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		klog.V(1).Infof("[%s] %s %s %v", req.RemoteAddr, req.Method, req.URL, req.Header)
		w.Write([]byte("ok"))
	})

	wrapped := std.Handler("default", mdlw, mux)

	server := &http.Server{
		Handler:      wrapped,
		ReadTimeout:  readTimeout,
		WriteTimeout: requestTimeout,
	}

	// Setup TLS if an hostname has been specified.
	if *tlsHostname != "" {
		// Proceed only if a valid hostname has been specified.
		if util.IsValidHostname(*tlsHostname) {
			klog.Fatalf("tlsHostname [%s] is an invalid hostname, exiting", *tlsHostname)
		}

		klog.Infof("provisioning Let's Encrypt certificate for %s", *tlsHostname)

		var acmeApi string
		if *tlsProd {
			klog.Infof("using production Let's Encrypt server")
			acmeApi = autocert.DefaultACMEDirectory
		} else {
			klog.Infof("using staging Let's Encrypt server")
			acmeApi = "https://acme-staging-v02.api.letsencrypt.org/directory"
		}

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(*tlsHostname),
			Cache:      autocert.DirCache(path.Join(*cacheDir, "autocert")),
			Client:     &acme.Client{DirectoryURL: acmeApi},
		}

		server.TLSConfig = certManager.TLSConfig()
		klog.Info("certificate provisioning configured")
	}

	// Graceful shutdown
	done := shutdownHandler(server)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		klog.Infof("debug server listening on %s", *debugAddr)
		klog.Exit(http.ListenAndServe(*debugAddr, nil))
	}()

	if *tlsHostname != "" {
		klog.Infof("main server listening with TLS on %s", listeners[0].Addr())
		klog.Exit(server.ServeTLS(listeners[0], "", ""))
	} else {
		klog.Infof("main server listening on %s", listeners[0].Addr())
		klog.Exit(server.Serve(listeners[0]))
	}

	<-done
}
