package routes

import (
	"ccp/backend/controllers"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("/ws", controllers.WebSocketHandler)
	mux.HandleFunc("/api/graph", controllers.GetElementsGraph)
	mux.HandleFunc("/api/elements", controllers.ElementsGetAll)

	// Serve static assets from "public"
	publicServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", publicServer))

	// Fallback: Serve frontend (SPA)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve static file from ./dist
		path := filepath.Join("./dist", r.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		// Fallback to index.html for client-side routing
		http.ServeFile(w, r, "./dist/index.html")
	})
}
