package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	//main api config
	port := "8080"
	cfg := &apiConfig{}
	cfg.fileServerHits.Store(0)

	//instanciate server mux
	serveMux := http.NewServeMux()

	// local handlers
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	//routing
	serveMux.Handle("/app/", cfg.middlewareMetricInc(fileServerHandler))
	serveMux.HandleFunc("/healthz", readinessHandler)
	serveMux.HandleFunc("/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("/reset", cfg.resetHandler)

	//server configuration
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	//server start
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileServerHits.Load())))
}
