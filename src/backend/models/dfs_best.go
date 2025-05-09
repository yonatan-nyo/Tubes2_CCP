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

	visited := map[string]bool{}
	tree, _, err := dfsBestTreeHelper(startNode, signalTreeChange, visited)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func dfsBestTreeHelper(
	targetNode *ElementsGraphNode,
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
	visited map[string]bool,
) (*RecipeTreeNode, int, error) {
	if IsBaseElement(targetNode.Name) {
		return &RecipeTreeNode{
			Name:                   targetNode.Name,
			ImagePath:              GetImagePath(targetNode.Name),
			MinimumNodesRecipeTree: 1,
		}, 1, nil
	}

	if len(targetNode.RecipesToMakeThisElement) == 0 {
		return nil, 0, fmt.Errorf("no recipes found to make element %s", targetNode.Name)
	}

	if visited[targetNode.Name] {
		return nil, 0, fmt.Errorf("cycle detected in the graph for element %s", targetNode.Name)
	}

	visited[targetNode.Name] = true
	defer delete(visited, targetNode.Name)

	var bestTree *RecipeTreeNode
	bestCost := int(^uint(0) >> 1)
	validPathFound := false

	for _, recipe := range targetNode.RecipesToMakeThisElement {
		// Avoid cycles
		if visited[recipe.ElementOne.Name] || visited[recipe.ElementTwo.Name] {
			continue
		}

		// Recursive search for both elements
		leftVisited := copyVisitedMap(visited)
		left, leftCost, errLeft := dfsBestTreeHelper(recipe.ElementOne, signalTreeChange, leftVisited)
		if errLeft != nil {
			continue
		}

		rightVisited := copyVisitedMap(visited)
		right, rightCost, errRight := dfsBestTreeHelper(recipe.ElementTwo, signalTreeChange, rightVisited)
		if errRight != nil {
			continue
		}

		validPathFound = true
		currentTree := &RecipeTreeNode{
			Name:                   targetNode.Name,
			ImagePath:              GetImagePath(targetNode.Name),
			Element1:               left,
			Element2:               right,
			MinimumNodesRecipeTree: leftCost + rightCost + 1,
		}

		if currentTree.MinimumNodesRecipeTree < bestCost {
			bestTree = currentTree
			bestCost = currentTree.MinimumNodesRecipeTree
			signalTreeChange(bestTree, currentTree)
		}
	}

	if !validPathFound {
		return nil, 0, fmt.Errorf("no valid recipe tree found for element %s", targetNode.Name)
	}

	return bestTree, bestCost, nil
}

func copyVisitedMap(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool, len(original))
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}
