package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OmarEP/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABSE_URL")
	if dbURL == "" {
		log.Fatal("PORT environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't load database:", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{ 
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth((apiCfg.handlerUsersGet)))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeedsGet)
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	mux.HandleFunc("GET /v1/posts", apiCfg.middlewareAuth(apiCfg.handlerPostGet))

	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))


	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	srv := &http.Server{
		Addr: 	":" + port,
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute 
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}