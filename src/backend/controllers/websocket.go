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

type FinalResponse struct {
	Trees         []*models.RecipeTreeNode `json:"trees"`
	DurationMs    int                      `json:"duration_ms"`
	NodesExplored int32                    `json:"nodes_explored"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

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

		updateChan := make(chan TreeUpdate, 1000)
		var latestUpdate *TreeUpdate
		var updateMu sync.Mutex

		signallerFn := func(
			exploringTree *models.RecipeTreeNode,
			durationMs int,
			nodesExplored int32,
		) {
			if req.DelayMs > 0 {
				updateMu.Lock()
				latestUpdate = &TreeUpdate{
					ExploringTree: exploringTree,
					DurationMs:    durationMs,
					NodesExplored: nodesExplored,
				}
				updateMu.Unlock()

				select {
				case updateChan <- *latestUpdate:
				default:
				}
			}
		}

		var updateWg sync.WaitGroup

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
		globalStartTime := time.Now()
		globalNodeCount := int32(0)
		trees, err := models.GenerateRecipeTree(req.Target, req.Mode, req.MaxTreeCount, signallerFn, req.DelayMs, globalStartTime, &globalNodeCount)

		close(updateChan)
		updateWg.Wait()

		writeMu.Lock()
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			writeMu.Unlock()
			continue
		}

		if err := conn.WriteJSON(
			FinalResponse{
				Trees:         trees,
				DurationMs:    int(time.Since(globalStartTime).Milliseconds()),
				NodesExplored: globalNodeCount,
			},
		); err != nil {
			log.Println("Final write error:", err)
			writeMu.Unlock()
			break
		}
		writeMu.Unlock()
	}
}
