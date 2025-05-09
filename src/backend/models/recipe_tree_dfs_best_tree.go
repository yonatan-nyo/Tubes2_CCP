package models

import "fmt"

func GenerateDFSFindBestTree(
	target string,
	signalTreeChange func(bestTree *RecipeTreeNode, exploringTree *RecipeTreeNode),
) (*RecipeTreeNode, error) {
	// Validate target node existence
	rootNode, ok := nameToNode[target]
	if !ok {
		return nil, fmt.Errorf("element %s not found in graph", target)
	}

	// Check if the target is a base element
	if IsBaseElement(target) {
		return &RecipeTreeNode{
			Name:                   target,
			ImagePath:              GetImagePath(target),
			MinimumNodesRecipeTree: 1,
		}, nil
	}

	// Setup map for generated trees
	computedTreeNode := make(map[string]*RecipeTreeNode)

	// Setup map for visited nodes to avoid cycles
	visited := make(map[string]bool)

	bestTree := &RecipeTreeNode{
		/* Initiate with the largest possible value to ensure any 
		   valid tree with smaller nodes will be chosen */
		MinimumNodesRecipeTree: int(^uint(0) >> 1), // Max int value
	}

	rootTree := &RecipeTreeNode{}

	// DFS process
	_, err := generateDFSFindBestTree(rootTree, rootNode, signalTreeChange, computedTreeNode, bestTree, visited)
	if err != nil {
		return nil, err
	}
	return rootTree, nil
}

// Helper function for DFS traversal to find the best tree
func generateDFSFindBestTree(
	currentRecipeTreeNode *RecipeTreeNode,
	currentGraphNode *ElementsGraphNode,
	signalTreeChange func(bestTree *RecipeTreeNode, exploringTree *RecipeTreeNode),
	computedTreeNode map[string]*RecipeTreeNode,
	bestTreeNode *RecipeTreeNode, // Track the best tree during traversal
	visited map[string]bool, // Track path-level visited nodes to avoid cycles
) (*RecipeTreeNode, error) {
	// Check if the current node is base element
	if IsBaseElement(currentGraphNode.Name) {
		currentRecipeTreeNode.Name = currentGraphNode.Name
		currentRecipeTreeNode.ImagePath = GetImagePath(currentGraphNode.ImagePath)
		currentRecipeTreeNode.MinimumNodesRecipeTree = 1
		return currentRecipeTreeNode, nil
	}

	// Cycle detection using visited map
	if visited[currentGraphNode.Name] {
		return nil, fmt.Errorf("cycle detected at element: %s", currentGraphNode.Name)
	}

	visited[currentGraphNode.Name] = true
	defer delete(visited, currentGraphNode.Name) // Allow reuse in other paths

	// Check cache (have already found the best tree for this node)
	if cached, ok := computedTreeNode[currentGraphNode.Name]; ok {
		*currentRecipeTreeNode = *cloneTree(cached)
		return currentRecipeTreeNode, nil
	}

	var bestLocalTree *RecipeTreeNode
	bestNodeCount := int(^uint(0) >> 1) // Max int value

	// Iterate through all recipes to find the best tree
	for _, currentRecipe := range currentGraphNode.RecipesToMakeThisElement {
		left := &RecipeTreeNode{}
		right := &RecipeTreeNode{}
		// Recursive call for both ingredients
		leftTree, err1 := generateDFSFindBestTree(left, currentRecipe.ElementOne, signalTreeChange, computedTreeNode, bestTreeNode, visited)
		rightTree, err2 := generateDFSFindBestTree(right, currentRecipe.ElementTwo, signalTreeChange, computedTreeNode, bestTreeNode, visited)

		if err1 != nil || err2 != nil {
			continue
		}

		nodeCount := leftTree.MinimumNodesRecipeTree + rightTree.MinimumNodesRecipeTree + 1

		// Check if the current combination is better than the best found so far
		if nodeCount < bestNodeCount {
			bestNodeCount = nodeCount
			bestLocalTree = &RecipeTreeNode{
				Name:                   currentGraphNode.Name,
				ImagePath:              GetImagePath(currentGraphNode.ImagePath),
				Element1:               leftTree,
				Element2:               rightTree,
				MinimumNodesRecipeTree: nodeCount,
			}

			// Update best tree so far (global)
			if nodeCount < bestTreeNode.MinimumNodesRecipeTree {
				*bestTreeNode = *cloneTree(bestLocalTree)
			}

			// Send update to WebSocket
			if signalTreeChange != nil {
				signalTreeChange(cloneTree(bestTreeNode), cloneTree(bestLocalTree))
			}

		}
	}

	if bestLocalTree == nil {
		return nil, fmt.Errorf("no valid recipe found for %s", currentGraphNode.Name)
	}

	// Cache result
	*currentRecipeTreeNode = *bestLocalTree
	computedTreeNode[currentGraphNode.Name] = cloneTree(bestLocalTree)

	return currentRecipeTreeNode, nil
}

/* Clone the tree to avoid modifying the original during traversal
   This is a deep copy function to ensure the original tree remains unchanged */
func cloneTree(node *RecipeTreeNode) *RecipeTreeNode {
	if node == nil {
		return nil
	}
	return &RecipeTreeNode{
		Name:                   node.Name,
		ImagePath:              node.ImagePath,
		Element1:               cloneTree(node.Element1),
		Element2:               cloneTree(node.Element2),
		MinimumNodesRecipeTree: node.MinimumNodesRecipeTree,
	}
}
