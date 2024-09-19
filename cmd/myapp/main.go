package main

import (
	"ether/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	log.Printf("Trying to serve: %s\n", path)

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		log.Printf("File not found, serving index.html\n")
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Serving file: %s\n", path)
	http.ServeFile(w, r, path)
}

func main() {
	checkEnv()

	r := mux.NewRouter()
	// r.PathPrefix("/").Handler(getStaticPath())

	r.HandleFunc("/block/number", handlers.GetBlockHandler)
	r.HandleFunc("/block/transactions", handlers.GetTransactionsHandler)
	r.HandleFunc("/balance", handlers.GetBalanceHandler)
	r.HandleFunc("/generate-key", handlers.GetAndGenerateKey)
	// r.HandleFunc("/balance/{address}/", handlers.GetBalanceHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + getPort(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on port %s in %s environment\n", getPort(), getEnvironment())
	log.Fatal(srv.ListenAndServe())
}

func checkEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}
}

func getEnvironment() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "production"
	}
	return environment
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func getStaticPath() spaHandler {
	staticPath := "static/dist"
	if getEnvironment() == "development" {
		staticPath = "static/dist"
	}

	spa := spaHandler{staticPath: staticPath, indexPath: "index.html"}
	return spa
}
