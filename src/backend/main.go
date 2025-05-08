package main

import (
	"ccp/backend/models"
	"ccp/backend/routes"
	"fmt"
	"log"
	"net/http"
)

// CORS middleware
func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4001")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	models.Init()
	mux := http.NewServeMux()

	models.Debug(models.ElementsGraph, -1)
	routes.RegisterRoutes(mux)

	// Wrap all routes with CORS
	handlerWithCORS := withCORS(mux)

	fmt.Println("Server started on :4000")
	log.Fatal(http.ListenAndServe("127.0.0.1:4000", handlerWithCORS))
}
