package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"strconv"
	"time"
)

type (
	rpcRequest struct {
		Version string        `json:"jsonrpc"`
		ID      interface{}   `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}
)

const (
	invalidRequestMessage = "Not a valid JSONRPC POST request\n"
	maxRequestSize        = 1024 * 10 // 10 KiB

	// Duration the entire request is allowed to take, including the backend request and response writing.
	requestTimeout = 10 * time.Second
	// Read timeout for request body and headers.
	readTimeout = 2 * time.Second
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "solrays_forwarded_duration_seconds",
		Help: "Duration of HTTP requests that made it to the backend",
	}, []string{"method"})

	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solrays_requests_total",
		Help: "Number of requests per method",
	}, []string{"method"})

	requestErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solrays_request_errors_total",
		Help: "Number of failed requests",
	}, []string{"reason"})

	backendRequest = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "solrays_requests_status_total",
		Help: "Number of backend responses by method and status code",
	}, []string{"method", "code"})
)

type handler struct {
	client *http.Client
}

func newHandler() handler {
	return handler{&http.Client{}}
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	klog.V(2).Infof("[%s] %s %s %v", req.RemoteAddr, req.Method, req.URL, req.Header)

	if req.Method == "OPTIONS" {
		h.sendCORS(w)
		return
	}

	if req.Method != "POST" {
		requestErrors.WithLabelValues("wrong_method").Inc()
		http.Error(w, invalidRequestMessage, http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		requestErrors.WithLabelValues("wrong_content_type").Inc()
		http.Error(w, invalidRequestMessage, http.StatusUnsupportedMediaType)
		return
	}

	h.sendCORS(w)

	ctx, cancel := context.WithTimeout(req.Context(), requestTimeout)
	defer cancel()

	// We need to read the entire request in order to decode and parse it.
	b, err := ioutil.ReadAll(http.MaxBytesReader(w, req.Body, maxRequestSize))
	if err != nil {
		requestErrors.WithLabelValues("reading_body").Inc()
		http.Error(w, "Error reading request body", http.StatusRequestTimeout)
		klog.V(1).Infof("[%s] (+%v) failed to read request body: %v",
			req.RemoteAddr, time.Since(start), err)
		return
	}

	klog.V(3).Infof("[%s] raw request: %s", req.RemoteAddr, string(b))

	var rpc rpcRequest
	if err = json.Unmarshal(b, &rpc); err != nil {
		requestErrors.WithLabelValues("decoding_request").Inc()
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		klog.V(1).Infof("[%s] (+%v) failed to decode request: %v",
			req.RemoteAddr, time.Since(start), err)
		return
	}

	klog.V(2).Infof("[%s] (+%v) request: %v", req.RemoteAddr, time.Since(start), rpc)

	if _, ok := methodWhitelistMap[rpc.Method]; !ok {
		requestErrors.WithLabelValues("invalid_method").Inc()
		http.Error(w, "Invalid RPC method", http.StatusBadRequest)
		klog.V(1).Infof("[%s] (+%v) called a method not on the whitelist: %s",
			req.RemoteAddr, time.Since(start), rpc.Method)
		return
	}

	requestsTotal.WithLabelValues(rpc.Method).Inc()

	r, err := http.NewRequestWithContext(ctx, "POST", *backendAddr, bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	r.Header = req.Header

	reqs := time.Now()
	res, err := h.client.Do(r)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		requestErrors.WithLabelValues("backend_request").Inc()
		klog.Warningf("[%s] (+%v) [%s] backend request failed: %v",
			req.RemoteAddr, rpc.ID, reqs.Sub(start), err)
		return
	}
	defer res.Body.Close()

	backendRequest.WithLabelValues(rpc.Method, strconv.Itoa(res.StatusCode)).Inc()
	klog.V(1).Infof("[%s] (+%v) [%s] %s request returned status %d in %v",
		req.RemoteAddr, time.Since(start), rpc.ID, rpc.Method, res.StatusCode, time.Since(reqs))
	w.WriteHeader(res.StatusCode)

	n, err := io.Copy(w, res.Body)
	klog.V(1).Infof("[%s] (+%v) [%s] wrote %d bytes", req.RemoteAddr, time.Since(start), rpc.ID, n)
	d := time.Since(start)
	if err != nil {
		klog.V(1).Infof("[%s] (+%v) [%s] failed writing response: %v",
			req.RemoteAddr, d, rpc.ID, err)
	}

	httpDuration.WithLabelValues(rpc.Method).Observe(d.Seconds())
}

func (h handler) sendCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
