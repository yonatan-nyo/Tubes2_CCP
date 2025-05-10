package models

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
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
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		return []*RecipeTreeNode{
			{
				Name:      targetGraphNode.Name,
				ImagePath: GetImagePath(targetGraphNode.ImagePath),
			},
		}, nil
	}

	queue := make(chan PartialTree, 1000) // Buffered channel to avoid blocking
	results := make(chan *RecipeTreeNode, maxTreeCount)
	var wg sync.WaitGroup

	var mu sync.Mutex
	resultCount := 0
	visitedTrees := make(map[string]bool) // Track visited tree configurations

	// Initial recipes
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		left := &RecipeTreeNode{Name: recipe.ElementOne.Name, ImagePath: GetImagePath(recipe.ElementOne.ImagePath)}
		right := &RecipeTreeNode{Name: recipe.ElementTwo.Name, ImagePath: GetImagePath(recipe.ElementTwo.ImagePath)}
		root := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
			Element1:  left,
			Element2:  right,
		}
		queue <- PartialTree{
			Node:    root,
			Pending: []*RecipeTreeNode{left, right},
		}
	}

	// Worker pool
	const workerCount = 8
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for current := range queue {
				mu.Lock()
				if resultCount >= maxTreeCount {
					mu.Unlock()
					return
				}
				mu.Unlock()

				pending := current.Pending
				if len(pending) == 0 {
					mu.Lock()
					if resultCount < maxTreeCount {
						results <- current.Node
						resultCount++
					}
					mu.Unlock()
					continue
				}

				toExpand := pending[0]
				remaining := pending[1:]

				graphNode := getElementByName(targetGraphNode, toExpand.Name)
				if graphNode == nil || len(graphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(graphNode.Name) {
					queue <- PartialTree{Node: current.Node, Pending: remaining}
					continue
				}

				// Check for new variations before processing
				for _, recipe := range graphNode.RecipesToMakeThisElement {
					left := &RecipeTreeNode{Name: recipe.ElementOne.Name, ImagePath: GetImagePath(recipe.ElementOne.ImagePath)}
					right := &RecipeTreeNode{Name: recipe.ElementTwo.Name, ImagePath: GetImagePath(recipe.ElementTwo.ImagePath)}

					newTree := current.Node.clone()
					newExpand := findNodeByName(newTree, toExpand.Name)
					if newExpand != nil {
						newExpand.Element1 = left
						newExpand.Element2 = right
					}

					// Check if the new tree configuration has been visited
					treeKey := newTree.GetKey() // Assuming GetKey() generates a unique identifier for the tree structure
					mu.Lock()
					if !visitedTrees[treeKey] && resultCount < maxTreeCount {
						visitedTrees[treeKey] = true
						queue <- PartialTree{
							Node:    newTree,
							Pending: append([]*RecipeTreeNode{left, right}, remaining...),
						}
					}
					mu.Unlock()
				}
			}
		}()
	}

	// Close queue when done
	go func() {
		wg.Wait()
		close(results)
	}()

	var final []*RecipeTreeNode
	for tree := range results {
		final = append(final, tree)
		if len(final) >= maxTreeCount {
			break
		}
	}

	return final, nil
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

func (node *RecipeTreeNode) GetKey() string {
	// Create a slice to hold the parts of the key (name and children's names)
	var parts []string
	parts = append(parts, node.Name)

	if node.Element1 != nil {
		parts = append(parts, node.Element1.GetKey())
	}
	if node.Element2 != nil {
		parts = append(parts, node.Element2.GetKey())
	}

	// Concatenate the parts into a single string (you can customize the separator if needed)
	keyString := strings.Join(parts, "|")

	// Create a SHA256 hash of the key string (you can also use MD5 if you prefer, though it is less secure)
	hash := sha256.New()
	hash.Write([]byte(keyString))
	return fmt.Sprintf("%x", hash.Sum(nil)) // Return the hex-encoded hash
}
