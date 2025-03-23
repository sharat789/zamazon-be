package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func newProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing target URL %s: %v", target, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Wrap the original director with additional debugging logs.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Log the target host and URL path that will be requested.
		log.Printf("[DEBUG] Forwarding to: %s%s", req.URL.Host, req.URL.Path)
	}

	// Add a ModifyResponse function to handle CORS headers
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		return nil
	}
	return proxy
}

func main() {
	// Create reverse proxies for each microservice
	usersProxy := newProxy("http://localhost:3000")
	buyerProxy := newProxy("http://localhost:3002")
	defaultProxy := newProxy("http://localhost:3001")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("Incoming request: %s", r.URL.Path)
		switch {
		case strings.HasPrefix(r.URL.Path, "/users"):
			usersProxy.ServeHTTP(w, r)
		case strings.HasPrefix(r.URL.Path, "/buyer"):
			buyerProxy.ServeHTTP(w, r)
		default:
			defaultProxy.ServeHTTP(w, r)
		}
	})

	log.Println("API Gateway is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
