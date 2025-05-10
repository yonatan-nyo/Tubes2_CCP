package models

import (
	"fmt"
	"sync"
)

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

	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
		}
		return []*RecipeTreeNode{node}, nil
	}

	var (
		result []*RecipeTreeNode
		mu     sync.Mutex
		count  = 0
		wg     sync.WaitGroup
	)

	treeChan := make(chan *RecipeTreeNode, maxTreeCount)
	queue := make([]PartialTree, 0)

	// Initialize the BFS queue with root combinations
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		left := &RecipeTreeNode{Name: recipe.ElementOne.Name, ImagePath: GetImagePath(recipe.ElementOne.ImagePath)}
		right := &RecipeTreeNode{Name: recipe.ElementTwo.Name, ImagePath: GetImagePath(recipe.ElementTwo.ImagePath)}
		root := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
			Element1:  left,
			Element2:  right,
		}
		queue = append(queue, PartialTree{
			Node:    root,
			Pending: []*RecipeTreeNode{left, right},
		})
	}

	// Process the queue using workers (goroutines)
	for len(queue) > 0 && count < maxTreeCount {
		nextQueue := []PartialTree{}

		for _, item := range queue {
			if len(item.Pending) == 0 {
				// Tree complete
				mu.Lock()
				if count < maxTreeCount {
					treeChan <- item.Node
					count++
				}
				mu.Unlock()
				continue
			}

			toExpand := item.Pending[0]
			rest := item.Pending[1:]

			graphNode := getElementByName(targetGraphNode, toExpand.Name)
			if graphNode == nil || IsBaseElement(graphNode.Name) || len(graphNode.RecipesToMakeThisElement) == 0 {
				nextQueue = append(nextQueue, PartialTree{
					Node:    item.Node,
					Pending: rest,
				})
				continue
			}

			// Expand all combinations from recipe
			for _, recipe := range graphNode.RecipesToMakeThisElement {
				wg.Add(1)
				go func(current PartialTree, rec *Recipe) {
					defer wg.Done()

					left := &RecipeTreeNode{Name: rec.ElementOne.Name, ImagePath: GetImagePath(rec.ElementOne.ImagePath)}
					right := &RecipeTreeNode{Name: rec.ElementTwo.Name, ImagePath: GetImagePath(rec.ElementTwo.ImagePath)}

					cloned := current.Node.clone()
					target := findNodeByName(cloned, toExpand.Name)
					if target != nil {
						target.Element1 = left
						target.Element2 = right
					}

					newPending := append([]*RecipeTreeNode{left, right}, rest...)
					mu.Lock()
					if count < maxTreeCount {
						nextQueue = append(nextQueue, PartialTree{
							Node:    cloned,
							Pending: newPending,
						})
					}
					mu.Unlock()
				}(item, recipe)
			}
		}

		wg.Wait()
		queue = nextQueue
	}

	close(treeChan)

	// Collect all trees from channel
	for tree := range treeChan {
		result = append(result, tree)
		if len(result) >= maxTreeCount {
			break
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