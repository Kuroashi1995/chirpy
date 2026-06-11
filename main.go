package main

import (
	"Kuroashi1995/chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	//env load
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	//database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database")
	}
	//main api config
	port := "8080"
	cfg := &apiConfig{}
	cfg.dbQueries = database.New(db)
	cfg.fileServerHits.Store(0)

	//instanciate server mux
	serveMux := http.NewServeMux()

	// local handlers
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	//routing
	//app
	serveMux.Handle("/app/", cfg.middlewareMetricInc(fileServerHandler))
	//admin
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	//api
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	//api - chirps
	serveMux.HandleFunc("POST /api/chirps", cfg.validateAndSaveChirp)
	serveMux.HandleFunc("GET /api/chirps/", cfg.getAllChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{id}", cfg.getChirpHandler)
	//api - user
	serveMux.HandleFunc("POST /api/users", cfg.registerUser)
	//api - auth
	serveMux.HandleFunc("POST /api/login", cfg.loginHandler)

	//server configuration
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	//server start
	log.Fatal(server.ListenAndServe())
}
