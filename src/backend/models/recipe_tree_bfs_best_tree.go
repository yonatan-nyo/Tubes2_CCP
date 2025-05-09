package models

import "fmt"

// *RecipeTreeNode is implemented in dfs file*
type RecipeTreeNodeWithCost struct {
	Tree  *RecipeTreeNode
	Cost  int
}

func GenerateBFSFindBestTree(
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

	// BFS process
	bestTree, err := generateBFSFindBestTree(rootNode, signalTreeChange)
	if err != nil {
		return nil, err
	}

	return bestTree, nil
}

// Helper function for BFS traversal to find the best tree
func generateBFSFindBestTree(
	targetNode *ElementsGraphNode,
	signalTreeChange func(bestTree *RecipeTreeNode, exploringTree *RecipeTreeNode),
) (*RecipeTreeNode, error) {

	// BFSNode holds the recipe tree at a specific node along with its total cost (node count)
	type BFSNode struct {
		Node *ElementsGraphNode
		Tree *RecipeTreeNode
		Cost int
	}

	// Initialize queue for BFS traversal and a map to track visited nodes with their minimum cost
	queue := []BFSNode{}
	visited := make(map[string]int)

	var bestTree *RecipeTreeNode
	bestCost := int(^uint(0) >> 1) // Initialize with max int value

	// Initialization: start the queue with all direct recipes to create the target element
	for _, recipe := range targetNode.RecipesToMakeThisElement {
		left := &RecipeTreeNode{
			Name:                   recipe.ElementOne.Name,
			ImagePath:              GetImagePath(recipe.ElementOne.ImagePath),
			MinimumNodesRecipeTree: 1,
		}
		right := &RecipeTreeNode{
			Name:                   recipe.ElementTwo.Name,
			ImagePath:              GetImagePath(recipe.ElementTwo.ImagePath),
			MinimumNodesRecipeTree: 1,
		}

		tree := &RecipeTreeNode{
			Name:                   targetNode.Name,
			ImagePath:              GetImagePath(targetNode.ImagePath),
			Element1:               left,
			Element2:               right,
			MinimumNodesRecipeTree: 3,
		}

		// Enqueue this initial tree for BFS
		queue = append(queue, BFSNode{
			Node: targetNode,
			Tree: tree,
			Cost: 3,
		})

		// Send initial exploration state to client (optional)
		if signalTreeChange != nil {
			signalTreeChange(nil, cloneTree(tree))
		}
	}

	// Begin BFS traversal
	for len(queue) > 0 {
		curr := queue[0] // Gets current node from the front of the queue
		queue = queue[1:] // Removes the current node from the queue

		// Skip if we have already found a cheaper way to build this element
		if cost, ok := visited[curr.Tree.Name]; ok && curr.Cost >= cost {
			continue
		}
		visited[curr.Tree.Name] = curr.Cost

		// If this is a better tree for the target element, update bestTree
		if curr.Cost < bestCost && curr.Tree.Name == targetNode.Name {
			bestTree = cloneTree(curr.Tree)
			bestCost = curr.Cost
			if signalTreeChange != nil {
				signalTreeChange(cloneTree(bestTree), cloneTree(curr.Tree))
			}
		}

		// Expand node, adding all recipes to make this element
		for _, recipe := range nameToNode[curr.Tree.Name].RecipesToMakeThisElement {
			leftTree := &RecipeTreeNode{
				Name:                   recipe.ElementOne.Name,
				ImagePath:              GetImagePath(recipe.ElementOne.ImagePath),
				MinimumNodesRecipeTree: 1,
			}
			rightTree := &RecipeTreeNode{
				Name:                   recipe.ElementTwo.Name,
				ImagePath:              GetImagePath(recipe.ElementTwo.ImagePath),
				MinimumNodesRecipeTree: 1,
			}
			newTree := &RecipeTreeNode{
				Name:                   curr.Tree.Name,
				ImagePath:              GetImagePath(nameToNode[curr.Tree.Name].ImagePath),
				Element1:               leftTree,
				Element2:               rightTree,
				MinimumNodesRecipeTree: leftTree.MinimumNodesRecipeTree + rightTree.MinimumNodesRecipeTree + 1,
			}
			queue = append(queue, BFSNode{
				Node: nameToNode[curr.Tree.Name],
				Tree: newTree,
				Cost: newTree.MinimumNodesRecipeTree,
			})

			// Send update for each exploration step
			if signalTreeChange != nil {
				signalTreeChange(cloneTree(bestTree), cloneTree(newTree))
			}
		}
	}

	if bestTree == nil {
		return nil, fmt.Errorf("no valid tree found for %s", targetNode.Name)
	}

	return bestTree, nil
}


