package controllers

import (
	"ccp/backend/models"
	"encoding/json"
	"net/http"
)

func RecipeTreeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	method := r.URL.Query().Get("method") // "dfs" or "bfs"

	if target == "" {
		http.Error(w, "Missing 'target' query param", http.StatusBadRequest)
		return
	}

	if method == "" {
		method = "dfs" // default
	}

	var (
		tree *models.RecipeTreeNode
		err  error
	)

	switch method {
	case "dfs":
		tree, err = models.GetRecipeTreeDFS(target)
	case "bfs":
		tree, err = models.GetRecipeTreeBFS(target)
	default:
		http.Error(w, "Unsupported method: "+method, http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to build recipe tree: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Type string                   `json:"type"`
		Tree *models.RecipeTreeNode   `json:"tree"`
	}{
		Type: "tree_" + method,
		Tree: tree,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
