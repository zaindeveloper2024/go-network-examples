package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type ProxyServer struct {
	client *http.Client
}

func NewProxyServer() *ProxyServer {
	return &ProxyServer{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dumpReq, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("Failed to dump request: %v", err)
	} else {
		log.Printf("Incoming Request:\n%s\n", string(dumpReq))
	}

	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create proxy request: %v", err), http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := proxyReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	resp, err := p.client.Do(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	dumpResp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("Failed to dump response: %v", err)
	} else {
		log.Printf("Outgoing Response:\n%s\n", string(dumpResp))
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to copy response: %v", err)
	}
}

func main() {
	proxy := NewProxyServer()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      proxy,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting proxy server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
