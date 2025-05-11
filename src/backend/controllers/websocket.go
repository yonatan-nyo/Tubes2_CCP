package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"
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
	NodesExplored int32                  `json:"nodes_explored"`
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

	// Add a mutex to protect writes to the websocket connection
	var writeMu sync.Mutex

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

		// Create a more reasonably sized buffer
		updateChan := make(chan TreeUpdate, 1000)
		// Track the latest update to ensure we don't miss any
		var latestUpdate *TreeUpdate
		var updateMu sync.Mutex

		// Signal function sends exploringTree to the channel
		signallerFn := func(
			exploringTree *models.RecipeTreeNode,
			durationMs int,
			nodesExplored int32,
		) {
			if req.DelayMs > 0 { // Only track updates if DelayMs is greater than 0
				// Store the latest update
				updateMu.Lock()
				latestUpdate = &TreeUpdate{
					ExploringTree: exploringTree,
					DurationMs:    durationMs,
					NodesExplored: nodesExplored,
				}
				updateMu.Unlock()

				// Also try to send to channel, but don't block if full
				select {
				case updateChan <- *latestUpdate:
				default:
					// If channel is full, that's ok, we'll send the latest update on the next tick
				}
			}
		}

		var updateWg sync.WaitGroup // <-- Add this

		// Launch goroutine to stream updates at intervals
		if req.DelayMs > 0 {
			updateWg.Add(1)
			go func() {
				defer updateWg.Done()
				ticker := time.NewTicker(time.Duration(req.DelayMs) * time.Millisecond)
				defer ticker.Stop()

				for update := range updateChan {
					<-ticker.C
					writeMu.Lock()
					if err := conn.WriteJSON(update); err != nil {
						log.Println("Write error:", err)
						writeMu.Unlock()
						return
					}
					writeMu.Unlock()
				}
			}()
		}

		// Generate the tree with live updates
		trees, err := models.GenerateRecipeTree(req.Target, req.Mode, req.MaxTreeCount, signallerFn, req.DelayMs)

		// Signal that tree generation is complete
		close(updateChan) // No more updates will be sent
		updateWg.Wait()   // ⬅️ Wait for the update-sending goroutine to finish

		// Now safely send the final tree after all updates
		writeMu.Lock()
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			writeMu.Unlock()
			continue
		}

		if err := conn.WriteJSON(trees); err != nil {
			log.Println("Final write error:", err)
			writeMu.Unlock()
			break
		}
		writeMu.Unlock()

	}
}
