package routes

import (
	"ccp/backend/controllers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws", controllers.WebSocketHandler)
}
