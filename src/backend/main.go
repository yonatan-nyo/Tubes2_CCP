package main

import (
	"ccp/backend/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	routes.RegisterRoutes(mux)

	fmt.Println("Server started on :4000")
	log.Fatal(http.ListenAndServe("127.0.0.1:4000", mux))
}
