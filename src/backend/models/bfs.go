package models

import (
	"fmt"
)

// PartialTree represents a work-in-progress tree with nodes yet to be expanded.
type PartialTree struct {
	Node    *RecipeTreeNode
	Pending []*RecipeTreeNode
}

func BFSFindTrees(
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode),
) ([]*RecipeTreeNode, error) {
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// If it's a base element, return directly
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 {
		return []*RecipeTreeNode{
			{
				Name:      targetGraphNode.Name,
				ImagePath: targetGraphNode.ImagePath,
			},
		}, nil
	}

	type PartialTree struct {
		Node    *RecipeTreeNode
		Pending []*RecipeTreeNode
	}

	queue := []PartialTree{}

	// Start with all recipes for the target element
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		left := &RecipeTreeNode{
			Name:      recipe.ElementOne.Name,
			ImagePath: recipe.ElementOne.ImagePath,
		}
		right := &RecipeTreeNode{
			Name:      recipe.ElementTwo.Name,
			ImagePath: recipe.ElementTwo.ImagePath,
		}
		root := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: targetGraphNode.ImagePath,
			Element1:  left,
			Element2:  right,
		}
		queue = append(queue, PartialTree{
			Node:    root,
			Pending: []*RecipeTreeNode{left, right},
		})
	}

	var result []*RecipeTreeNode

	for len(queue) > 0 && len(result) < maxTreeCount {
		current := queue[0]
		queue = queue[1:]

		pending := current.Pending

		if len(pending) == 0 {
			// Fully constructed tree
			result = append(result, current.Node)
			continue
		}

		// Expand the first pending node
		toExpand := pending[0]
		remaining := pending[1:]

		graphNode := getElementByName(targetGraphNode, toExpand.Name)
		if graphNode == nil {
			continue // Skip unknown nodes
		}

		// If base element, mark and continue
		if len(graphNode.RecipesToMakeThisElement) == 0 {
			queue = append(queue, PartialTree{Node: current.Node, Pending: remaining})
			continue
		}

		// Try all recipes for this node
		for _, recipe := range graphNode.RecipesToMakeThisElement {
			// Create new subnodes
			left := &RecipeTreeNode{
				Name:      recipe.ElementOne.Name,
				ImagePath: recipe.ElementOne.ImagePath,
			}
			right := &RecipeTreeNode{
				Name:      recipe.ElementTwo.Name,
				ImagePath: recipe.ElementTwo.ImagePath,
			}

			// Clone full tree
			newTree := current.Node.clone()

			// Find the matching node to expand in the new tree
			newExpand := findNodeByName(newTree, toExpand.Name)
			if newExpand != nil {
				newExpand.Element1 = left
				newExpand.Element2 = right
			}

			newPending := append([]*RecipeTreeNode{left, right}, remaining...)
			queue = append(queue, PartialTree{
				Node:    newTree,
				Pending: newPending,
			})
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid trees found for %s", targetGraphNode.Name)
	}

	return result, nil
}

func findNodeByName(root *RecipeTreeNode, name string) *RecipeTreeNode {
	if root == nil {
		return nil
	}
	if root.Name == name {
		return root
	}
	if found := findNodeByName(root.Element1, name); found != nil {
		return found
	}
	return findNodeByName(root.Element2, name)
}

func getElementByName(root *ElementsGraphNode, name string) *ElementsGraphNode {
	if root == nil {
		return nil
	}
	if root.Name == name {
		return root
	}
	for _, recipe := range root.RecipesToMakeThisElement {
		if found := getElementByName(recipe.ElementOne, name); found != nil {
			return found
		}
		if found := getElementByName(recipe.ElementTwo, name); found != nil {
			return found
		}
	}
	return nil
}
