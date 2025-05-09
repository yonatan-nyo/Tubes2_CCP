package models

import (
	"fmt"
	"math"
)

func DFSFindTrees(
	rootRecipeTree *RecipeTreeNode,
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode),
) ([]*RecipeTreeNode, error) {
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// Primitive base elements (cannot be crafted further)
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 {
		node := &RecipeTreeNode{
			Name:                   targetGraphNode.Name,
			ImagePath:              targetGraphNode.ImagePath,
			MinimumNodesRecipeTree: 0,
		}
		return []*RecipeTreeNode{node}, nil
	}

	minNodeCount := math.MaxInt32
	var bestTrees []*RecipeTreeNode

	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		leftTrees, err1 := DFSFindTrees(nil, recipe.ElementOne, maxTreeCount, signalTreeChange)
		rightTrees, err2 := DFSFindTrees(nil, recipe.ElementTwo, maxTreeCount, signalTreeChange)

		if err1 != nil || err2 != nil {
			continue
		}

		for _, lt := range leftTrees {
			for _, rt := range rightTrees {
				root := &RecipeTreeNode{
					Name:                   targetGraphNode.Name,
					ImagePath:              targetGraphNode.ImagePath,
					Element1:               lt,
					Element2:               rt,
					MinimumNodesRecipeTree: lt.MinimumNodesRecipeTree + rt.MinimumNodesRecipeTree + 1,
				}

				if root.MinimumNodesRecipeTree < minNodeCount {
					minNodeCount = root.MinimumNodesRecipeTree
					bestTrees = []*RecipeTreeNode{root}
					// signalTreeChange(root)
				} else if root.MinimumNodesRecipeTree == minNodeCount {
					bestTrees = append(bestTrees, root)
				}

				// Limit results if maxTreeCount is reached
				if len(bestTrees) >= maxTreeCount {
					return bestTrees, nil
				}
			}
		}
	}

	return bestTrees, nil
}
