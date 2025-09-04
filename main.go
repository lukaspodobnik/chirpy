package main

// connection-string: "postgres://lukaspodobnik:@localhost:5432/chirpy"

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/lukaspodobnik/chirpy/internal/database"
)

const (
	filepathroot = "."
	port         = "8080"
)

const (
	contenttype   = "Content-Type"
	textPlainUTF8 = "text/plain; charset=utf-8"
	textHtml      = "text/html"
	appJson       = "application/json"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading environment failed: %v\n", err)
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("opening database failed: %v\n", err)
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
	}

	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathroot))))

	mux := http.NewServeMux()

	mux.Handle("/app/", fileserverHandler)

	mux.HandleFunc("GET /api/healthz", getHealthzHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.postUserHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.getMetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.postResetHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files from %s on port: %s\n", filepathroot, port)
	log.Fatal(server.ListenAndServe())
}
