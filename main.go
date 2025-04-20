package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var proxy *httputil.ReverseProxy

func main() {
	// Get target URL from environment variable or fallback
	target := os.Getenv("TARGET_URL")
	if target == "" {
		target = "http://yamabiko.proxy.rlwy.net:36864"
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid TARGET_URL: %v", err)
	}

	proxy = httputil.NewSingleHostReverseProxy(targetURL)

	// Route definitions
	http.HandleFunc("/send-email", withCORS(proxyHandler))
	http.HandleFunc("/send-single-email", withCORS(proxyHandler))
	http.HandleFunc("/add-recipe", withCORS(proxyHandler))
	http.HandleFunc("/get-all-recipes", withCORS(proxyHandler))
	http.HandleFunc("/generate-recipes", withCORS(proxyHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🔁 Secure Reverse Proxy server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// proxyHandler proxies the request to the target backend
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r)
}

// withCORS adds CORS headers
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}
