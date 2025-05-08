package models

import "fmt"

func GenerateDFSFindBestTree(
	target string,
	signalTreeChange func(bestTree *RecipeTreeNode, exploringTree *RecipeTreeNode),
) (*RecipeTreeNode, error) {
	rootRecipeTree := &RecipeTreeNode{
		Name:      target,
		ImagePath: GetImagePath(target),
	}
	if IsBaseElement(target) {
		rootRecipeTree.MinimumNodesRecipeTree = 1
		return rootRecipeTree, nil
	}
	return nil, nil
}

// Helper function for DFS traversal to find the best tree
func generateDFSFindBestTree(
	currentRecipeTreeNode *RecipeTreeNode,
	currentGraphNode *ElementsGraphNode,
	signalTreeChange func(bestTree *RecipeTreeNode, exploringTree *RecipeTreeNode),
	computedTreeNode map[string]*RecipeTreeNode,
	bestTreeNode *RecipeTreeNode, // Track the best tree during traversal
) (*RecipeTreeNode, error) {
	if currentGraphNode == nil {
		return nil, fmt.Errorf("currentGraphNode is nil")
	}
	if currentRecipeTreeNode == nil {
		return nil, fmt.Errorf("currentRecipeTreeNode is nil")
	}

	return nil, nil
}
