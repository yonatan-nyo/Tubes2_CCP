package models

import "errors"

// Struct RecipeTreeNode
type RecipeTreeNode struct {
	Name      string          `json:"name"`
	ImagePath string          `json:"image_path"`
	Element1  *RecipeTreeNode `json:"element_1,omitempty"`
	Element2  *RecipeTreeNode `json:"element_2,omitempty"`
}

// ErrElementNotFound: element not found in graph
func ErrElementNotFound(name string) error {
	return errors.New("element not found: " + name)
}

// Return tree with least node count
func GetRecipeTree(target string) (*RecipeTreeNode, error) {
	// Value map (validasi)
	targetNode, exists := nameToNode[target]
	if !exists {
		return nil, ErrElementNotFound(target)
	}

	// Base element: Air, Water, Earth, Fire
	if len(targetNode.RecipesToMakeThisElement) == 0 {
		return &RecipeTreeNode{
			Name:      targetNode.Name,
			ImagePath: GetImagePath(targetNode.ImagePath),
		}, nil
	}

	var bestTree *RecipeTreeNode
	minNodeCount := -1

	for _, recipe := range targetNode.RecipesToMakeThisElement {
		visited := map[string]bool{} // Avoid infinite cycle
		tree := &RecipeTreeNode{
			Name:      targetNode.Name,
			ImagePath: GetImagePath(targetNode.ImagePath),
			Element1:  buildRecipeTree(recipe.ElementOne, visited),
			Element2:  buildRecipeTree(recipe.ElementTwo, visited),
		}

		count := countRecipeTreeNode(tree)

		if minNodeCount == -1 || count < minNodeCount {
			minNodeCount = count
			bestTree = tree
		}
	}

	if bestTree == nil {
		return nil, errors.New("no valid recipe tree could be constructed")
	}

	return bestTree, nil
}

func buildRecipeTree(node *ElementsGraphNode, visited map[string]bool) *RecipeTreeNode {
	// Basis 1: node kosong
	if node == nil {
		return nil
	}

	// Always create nodes for base elements regardless of visited status
	if isBaseElement(node.Name) {
		return &RecipeTreeNode{
			Name:      node.Name,
			ImagePath: GetImagePath(node.ImagePath),
		}
	}

	// Untuk non-base elements
	if visited[node.Name] {
		return nil
	}
	visited[node.Name] = true

	// Base element or no recipe
	if len(node.RecipesToMakeThisElement) == 0 {
		return &RecipeTreeNode{
			Name:      node.Name,
			ImagePath: GetImagePath(node.ImagePath),
		}
	}

	// Use the first recipe
	best := node.RecipesToMakeThisElement[0]

	tree := &RecipeTreeNode{
		Name:      node.Name,
		ImagePath: GetImagePath(node.ImagePath),
	}

	// Create a copy of the visited map for each branch to prevent cross-branch interference
	visitedElement1 := make(map[string]bool)
	visitedElement2 := make(map[string]bool)

	for k, v := range visited {
		visitedElement1[k] = v
		visitedElement2[k] = v
	}

	// Proses element pertama
	if best.ElementOne != nil {
		tree.Element1 = buildRecipeTree(best.ElementOne, visitedElement1)
	}

	// Proses element kedua
	if best.ElementTwo != nil {
		tree.Element2 = buildRecipeTree(best.ElementTwo, visitedElement2)
	}

	return tree
}

// Count Recipe Tree node
func countRecipeTreeNode(tree *RecipeTreeNode) int {
	if tree == nil {
		return 0
	}
	return 1 + countRecipeTreeNode(tree.Element1) + countRecipeTreeNode(tree.Element2)
}

// Check if the element is a base element
// Base elements: Air, Water, Earth, Fire
func isBaseElement(name string) bool {
	return name == "Air" || name == "Water" || name == "Earth" || name == "Fire"
}
