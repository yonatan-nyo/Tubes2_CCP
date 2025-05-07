package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"net/http"
)

func GetElementsGraph(w http.ResponseWriter, r *http.Request) {
	if models.ElementsGraph == nil {
		http.Error(w, "Graph not found", http.StatusNotFound)
		return
	}

	safeGraphNode := models.GetJSONDTONodes()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(safeGraphNode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
