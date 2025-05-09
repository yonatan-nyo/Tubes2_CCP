package models

import "fmt"

func GenerateDFSFindBestTree(
	target string,
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
) (*RecipeTreeNode, error) {
	startNode, ok := nameToNode[target]
	if !ok || startNode == nil {
		return nil, fmt.Errorf("target %s not found in elements graph", target)
	}

	parrentElementNameMap := make(map[string]bool)
	tree, _, err := dfsBestTreeHelper(startNode, signalTreeChange, parrentElementNameMap, target)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func dfsBestTreeHelper(
	targetNode *ElementsGraphNode,
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
	parentParentElement map[string]bool,
	parentElementName string,
) (*RecipeTreeNode, int, error) {
	if IsBaseElement(targetNode.Name) {
		return &RecipeTreeNode{
			Name:                   targetNode.Name,
			ImagePath:              GetImagePath(targetNode.Name),
			MinimumNodesRecipeTree: 1,
		}, 1, nil
	}

	// check for recipes to make this element
	if len(targetNode.RecipesToMakeThisElement) == 0 {
		return nil, 0, fmt.Errorf("no recipes found to make element %s", targetNode.Name)
	}

	// Check if the target node has already been visited
	if _, ok := parentParentElement[targetNode.Name]; ok {
		return nil, 0, fmt.Errorf("cycle detected in the graph for element %s", targetNode.Name)
	}

	// Mark the current node as visited
	parentParentElement[parentElementName] = true
	defer func() {
		// Unmark the current node when backtracking
		delete(parentParentElement, parentElementName)
	}()

	// Initialize the best tree with the first recipe
	bestMinimumNodesRecipeTreePlaceholder := int(^uint(0) >> 1) // Initialize with max int value
	bestTree := &RecipeTreeNode{
		Name:                   targetNode.Name,
		ImagePath:              GetImagePath(targetNode.Name),
		MinimumNodesRecipeTree: bestMinimumNodesRecipeTreePlaceholder,
	}
	bestTree.Element1 = &RecipeTreeNode{
		Name:                   targetNode.RecipesToMakeThisElement[0].ElementOne.Name,
		ImagePath:              GetImagePath(targetNode.RecipesToMakeThisElement[0].ElementOne.ImagePath),
		MinimumNodesRecipeTree: bestMinimumNodesRecipeTreePlaceholder,
	}
	bestTree.Element2 = &RecipeTreeNode{
		Name:                   targetNode.RecipesToMakeThisElement[0].ElementTwo.Name,
		ImagePath:              GetImagePath(targetNode.RecipesToMakeThisElement[0].ElementTwo.ImagePath),
		MinimumNodesRecipeTree: bestMinimumNodesRecipeTreePlaceholder,
	}
	// find the element1 and 2 actual cost with recursive call
	left, leftCost, err := dfsBestTreeHelper(
		targetNode.RecipesToMakeThisElement[0].ElementOne,
		signalTreeChange,
		parentParentElement,
		targetNode.Name,
	)
	if err != nil {
		return nil, 0, err
	}
	right, rightCost, err := dfsBestTreeHelper(
		targetNode.RecipesToMakeThisElement[0].ElementTwo,
		signalTreeChange,
		parentParentElement,
		targetNode.Name,
	)
	if err != nil {
		return nil, 0, err
	}
	bestTree.Element1 = left
	bestTree.Element2 = right
	bestTree.MinimumNodesRecipeTree = leftCost + rightCost + 1

	// loop through all recipes to find the best tree
	for _, recipe := range targetNode.RecipesToMakeThisElement[1:] {
		left, leftCost, err := dfsBestTreeHelper(
			recipe.ElementOne,
			signalTreeChange,
			parentParentElement,
			targetNode.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		right, rightCost, err := dfsBestTreeHelper(
			recipe.ElementTwo,
			signalTreeChange,
			parentParentElement,
			targetNode.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		currentTree := &RecipeTreeNode{
			Name:                   targetNode.Name,
			ImagePath:              GetImagePath(targetNode.Name),
			Element1:               left,
			Element2:               right,
			MinimumNodesRecipeTree: leftCost + rightCost + 1,
		}
		// Check if the current tree is better than the best tree
		if currentTree.MinimumNodesRecipeTree < bestTree.MinimumNodesRecipeTree {
			bestTree = currentTree
			bestMinimumNodesRecipeTreePlaceholder = currentTree.MinimumNodesRecipeTree
			signalTreeChange(bestTree, currentTree)
		}
	}

	return bestTree, bestMinimumNodesRecipeTreePlaceholder, nil
}
