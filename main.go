package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits counter reset."))
}

func main() {
	mux := http.NewServeMux()
	cfg := apiConfig{}

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("/reset", cfg.resetHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
