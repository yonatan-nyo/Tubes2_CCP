package models

import (
	"fmt"
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

	// Base case: base element
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: targetGraphNode.ImagePath,
		}
		return []*RecipeTreeNode{node}, nil
	}

	var allTrees []*RecipeTreeNode

	// Explore all recipes for this element
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		if len(allTrees) >= maxTreeCount {
			break
		}

		leftTrees, err1 := DFSFindTrees(nil, recipe.ElementOne, maxTreeCount, signalTreeChange)
		if err1 != nil {
			continue
		}

		rightTrees, err2 := DFSFindTrees(nil, recipe.ElementTwo, maxTreeCount, signalTreeChange)
		if err2 != nil {
			continue
		}

		for _, lt := range leftTrees {
			if len(allTrees) >= maxTreeCount {
				break
			}
			for _, rt := range rightTrees {
				if len(allTrees) >= maxTreeCount {
					break
				}

				root := &RecipeTreeNode{
					Name:      targetGraphNode.Name,
					ImagePath: targetGraphNode.ImagePath,
					Element1:  lt,
					Element2:  rt,
				}

				allTrees = append(allTrees, root)
			}
		}
	}

	if len(allTrees) == 0 {
		return nil, fmt.Errorf("no valid trees found for %s", targetGraphNode.Name)
	}

	return allTrees, nil
}
