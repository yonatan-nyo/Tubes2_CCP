package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

type RecipeTreeRequest struct {
	Target       string `json:"target"`
	Mode         string `json:"mode"`
	MaxTreeCount int    `json:"max_tree_count"`
	DelayMs      int    `json:"delay_ms"`
}

type TreeUpdate struct {
	ExploringTree *models.RecipeTreeNode `json:"exploring_tree"`
	DurationMs    int                    `json:"duration_ms"`
	NodesExplored int32 				 `json:"nodes_explored"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var req RecipeTreeRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Println("Invalid JSON format:", err)
			continue
		}

		// Create channel to buffer recipe tree updates
		updateChan := make(chan TreeUpdate, 1000000)

		// Signal function sends both bestTree and exploringTree to the channel
		signallerFn := func(
			exploringTree *models.RecipeTreeNode,
			durationMs int,
			nodesExplored int32,
		) {
			if req.DelayMs > 0 { // Only send updates if DelayMs is greater than 0
				select {
				case updateChan <- TreeUpdate{
					ExploringTree: exploringTree,
					DurationMs: durationMs,
					NodesExplored: nodesExplored,
					}:
				default:
					log.Println("Warning: updateChan full, dropping update")
				}
			}
		}

		// Launch goroutine to stream updates at intervals
		if req.DelayMs > 0 {
			go func() {
				ticker := time.NewTicker(time.Duration(req.DelayMs) * time.Millisecond)
				defer ticker.Stop()

				for update := range updateChan {
					<-ticker.C
					if err := conn.WriteJSON(update); err != nil {
						log.Println("Write error:", err)
						return
					}
				}
			}()
		}

		// Generate the tree with live updates
		trees, err := models.GenerateRecipeTree(req.Target, req.Mode, req.MaxTreeCount, signallerFn)
		close(updateChan) // close after final tree built

		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			continue
		}

		// Final full tree send (optional)
		if err := conn.WriteJSON(trees); err != nil {
			log.Println("Final write error:", err)
			break
		}
	}
}