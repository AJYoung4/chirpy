package main

import (
	"fmt"
	"net/http"

	"github.com/AJYoung4/chirpy/api"
)

func main() {
	mux := http.NewServeMux()
	cfg := api.ApiConfig{}

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fs)))

	mux.HandleFunc("GET /api/healthz", api.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
