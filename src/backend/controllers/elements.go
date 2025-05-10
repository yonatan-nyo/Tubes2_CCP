package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"net/http"
)

func ElementsGetAll(w http.ResponseWriter, r *http.Request) {
	if models.ElementsGraph == nil {
		http.Error(w, "Graph not found", http.StatusNotFound)
		return
	}

	elements := models.GetElementsFromNameToNodeDTO()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(elements); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
