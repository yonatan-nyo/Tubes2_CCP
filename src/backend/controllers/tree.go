package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"net/http"
)

// RecipeTreeHandler: retrieve Recipe Tree with the least node count
func RecipeTreeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "Missing 'target' query param", http.StatusBadRequest)
		return
	}

	tree, err := models.GetRecipeTree(target)
	if err != nil {
		http.Error(w, "Failed to build recipe tree: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Type string                 `json:"type"`
		Tree *models.RecipeTreeNode `json:"tree"`
	}{
		Type: "tree",
		Tree: tree,
	}

	if tree == nil {
		http.Error(w, "No tree could be constructed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
