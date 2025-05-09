package models

import "fmt"

var logCounter = 0

// dfsBuildTree is a helper that recursively constructs the best recipe tree using DFS.
func GenerateDFSFindBestTree(
	target string,
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
) (*RecipeTreeNode, error) {
	startNode, ok := nameToNode[target]
	if !ok || startNode == nil {
		return nil, fmt.Errorf("target %s not found in elements graph", target)
	}

	visiting := make(map[string]bool)
	tree := dfsBuildTree(startNode, visiting, signalTreeChange)

	if tree == nil {
		return nil, fmt.Errorf("no valid recipe tree found for %s", target)
	}

	return tree, nil
}

func dfsBuildTree(
	node *ElementsGraphNode,
	visiting map[string]bool,
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode),
) *RecipeTreeNode {
	if node == nil {
		return nil
	}

	if logCounter < 50 {
		fmt.Println("building tree for:", node.Name)
	}

	if visiting[node.Name] {
		if logCounter < 50 {
			fmt.Println("cycle detected at:", node.Name)
			logCounter++
		}
		return nil
	}
	visiting[node.Name] = true

	if len(node.RecipesToMakeThisElement) == 0 {
		visiting[node.Name] = false
		return &RecipeTreeNode{
			Name:                   node.Name,
			ImagePath:              node.ImagePath,
			MinimumNodesRecipeTree: 1,
			IsParentElement:        make(map[string]bool),
		}
	}

	minNodes := int(^uint(0) >> 1)
	var bestTree *RecipeTreeNode

	for _, recipe := range node.RecipesToMakeThisElement {
		if logCounter < 50 {
			fmt.Println(" trying recipe ", recipe.TargetElementName)
		}

		// PRE-CHECK for cycle before recursion
		if visiting[recipe.ElementOne.Name] || visiting[recipe.ElementTwo.Name] {
			continue // Skip recipe entirely if either causes a cycle
		}

		left := dfsBuildTree(recipe.ElementOne, visiting, signalTreeChange)
		right := dfsBuildTree(recipe.ElementTwo, visiting, signalTreeChange)

		if left != nil && right != nil {
			total := left.MinimumNodesRecipeTree + right.MinimumNodesRecipeTree + 1
			if total < minNodes {
				bestTree = &RecipeTreeNode{
					Name:                   node.Name,
					ImagePath:              node.ImagePath,
					Element1:               left,
					Element2:               right,
					MinimumNodesRecipeTree: total,
					IsParentElement: map[string]bool{
						left.Name:  true,
						right.Name: true,
					},
				}
				signalTreeChange(left, bestTree)
				signalTreeChange(right, bestTree)
				minNodes = total
			}
		}
	}

	visiting[node.Name] = false
	return bestTree
}
