package models

// func GenerateDFSFindBestTree(
// 	target string,
// 	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
// ) (*RecipeTreeNode, error) {
// 	startNode, ok := nameToNode[target]
// 	if !ok || startNode == nil {
// 		return nil, fmt.Errorf("target %s not found in elements graph", target)
// 	}

// 	currentCompleteTreeRoot := &RecipeTreeNode{
// 		Name:                   startNode.Name,
// 		ImagePath:              GetImagePath(startNode.Name),
// 		MinimumNodesRecipeTree: int(^uint(0) >> 1),
// 	}
// 	bestCompleteTreeRoot := currentCompleteTreeRoot.clone()

// 	signallerFn := func() {
// 		signalTreeChange(currentCompleteTreeRoot, bestCompleteTreeRoot)
// 	}

// 	visited := map[string]bool{}
// 	tree, _, err := dfsBestTreeHelper(startNode, currentCompleteTreeRoot, signallerFn, visited, currentCompleteTreeRoot, bestCompleteTreeRoot)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tree, nil
// }

// func dfsBestTreeHelper(
// 	targetGraphNode *ElementsGraphNode,
// 	targetRecipeNode *RecipeTreeNode,
// 	signalTreeChange func(),
// 	visited map[string]bool,
// 	currentCompleteTreeRoot *RecipeTreeNode,
// 	bestCompleteTreeRoot *RecipeTreeNode,
// ) (*RecipeTreeNode, int, error) {
// 	if targetGraphNode == nil {
// 		return nil, 0, fmt.Errorf("targetGraphNode is nil")
// 	}
// 	if targetRecipeNode == nil {
// 		return nil, 0, fmt.Errorf("targetRecipeNode is nil")
// 	}
// 	if IsBaseElement(targetGraphNode.Name) {
// 		return &RecipeTreeNode{
// 			Name:                   targetGraphNode.Name,
// 			ImagePath:              GetImagePath(targetGraphNode.Name),
// 			MinimumNodesRecipeTree: 1,
// 		}, 1, nil
// 	}

// 	if len(targetGraphNode.RecipesToMakeThisElement) == 0 {
// 		return nil, 0, fmt.Errorf("no recipes found to make element %s", targetGraphNode.Name)
// 	}

// 	if visited[targetGraphNode.Name] {
// 		return nil, 0, fmt.Errorf("cycle detected in the graph for element %s", targetGraphNode.Name)
// 	}

// 	visited[targetGraphNode.Name] = true
// 	defer delete(visited, targetGraphNode.Name)

// 	var bestTree *RecipeTreeNode
// 	bestCost := int(^uint(0) >> 1)
// 	validPathFound := false

// 	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
// 		if visited[recipe.ElementOne.Name] || visited[recipe.ElementTwo.Name] {
// 			continue
// 		}

// 		el1 := &RecipeTreeNode{
// 			Name:                   recipe.ElementOne.Name,
// 			ImagePath:              GetImagePath(recipe.ElementOne.Name),
// 			MinimumNodesRecipeTree: 1,
// 		}
// 		targetRecipeNode.Element1 = el1
// 		leftVisited := copyVisitedMap(visited)
// 		left, leftCost, err := dfsBestTreeHelper(recipe.ElementOne, el1, signalTreeChange, leftVisited, currentCompleteTreeRoot, bestCompleteTreeRoot)
// 		if err != nil {
// 			continue
// 		}

// 		el2 := &RecipeTreeNode{
// 			Name:                   recipe.ElementTwo.Name,
// 			ImagePath:              GetImagePath(recipe.ElementTwo.Name),
// 			MinimumNodesRecipeTree: 1,
// 		}
// 		targetRecipeNode.Element2 = el2
// 		rightVisited := copyVisitedMap(visited)
// 		right, rightCost, err := dfsBestTreeHelper(recipe.ElementTwo, el2, signalTreeChange, rightVisited, currentCompleteTreeRoot, bestCompleteTreeRoot)
// 		if err != nil {
// 			continue
// 		}

// 		validPathFound = true
// 		targetRecipeNode.Element1 = left
// 		targetRecipeNode.Element2 = right
// 		targetRecipeNode.MinimumNodesRecipeTree = leftCost + rightCost + 1

// 		if targetRecipeNode.MinimumNodesRecipeTree < bestCost {
// 			bestTree = targetRecipeNode
// 			bestCost = targetRecipeNode.MinimumNodesRecipeTree

// 			// Clone current full tree into best complete tree
// 			*bestCompleteTreeRoot = *currentCompleteTreeRoot.clone()

// 			signalTreeChange()
// 		}
// 	}
// 	if bestTree != nil {
// 		bestTree.Element1 = targetRecipeNode.Element1
// 		bestTree.Element2 = targetRecipeNode.Element2
// 		bestTree.MinimumNodesRecipeTree = bestCost
// 	} else {
// 		return nil, 0, fmt.Errorf("no valid recipe tree found for element %s", targetGraphNode.Name)
// 	}

// 	if !validPathFound {
// 		return nil, 0, fmt.Errorf("no valid recipe tree found for element %s", targetGraphNode.Name)
// 	}

// 	return bestTree, bestCost, nil
// }

// func copyVisitedMap(original map[string]bool) map[string]bool {
// 	newMap := make(map[string]bool, len(original))
// 	for k, v := range original {
// 		newMap[k] = v
// 	}
// 	return newMap
// }
