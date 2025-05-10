package routes

import (
	"ccp/backend/controllers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", controllers.WebSocketHandler)
	mux.HandleFunc("/api/graph", controllers.GetElementsGraph)
	mux.HandleFunc("/api/elements", controllers.ElementsGetAll)

	// Serve static files from the "public" folder
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))
}
