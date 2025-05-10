package models

import (
	"fmt"
	"sync"
)

// BFSFindTrees implements BFS approach to find recipe trees
func BFSFindTrees(
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode),
) ([]*RecipeTreeNode, error) {
	if targetGraphNode == nil {
		return nil, fmt.Errorf("targetGraphNode is nil")
	}

	// If the target is a base element or has no recipes, return it directly
	if len(targetGraphNode.RecipesToMakeThisElement) == 0 || IsBaseElement(targetGraphNode.Name) {
		node := &RecipeTreeNode{
			Name:      targetGraphNode.Name,
			ImagePath: GetImagePath(targetGraphNode.ImagePath),
		}
		return []*RecipeTreeNode{node}, nil
	}

	// Variables for tracking results and synchronization
	var (
		result     []*RecipeTreeNode
		mu         sync.Mutex
		treesFound int = 0
	)

	// Queue for BFS
	type QueueItem struct {
		ElementName string
		Level       int
		TreeSoFar   *RecipeTreeNode
		IsComplete  bool
	}

	// Process each target recipe concurrently
	var wg sync.WaitGroup
	resultChan := make(chan *RecipeTreeNode, maxTreeCount)

	// Process each initial recipe for the target element
	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			// Map to store all possible trees for each element
			elementToTrees := make(map[string][]*RecipeTreeNode)
			processedElements := make(map[string]bool)
			
			// Queue for BFS traversal
			queue := make([]*QueueItem, 0)
			
			// Start with the elements needed for this recipe
			queue = append(
				queue,
				&QueueItem{ElementName: r.ElementOne.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
				&QueueItem{ElementName: r.ElementTwo.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
			)
			
			// BFS traversal
			for len(queue) > 0 {
				mu.Lock()
				if treesFound >= maxTreeCount {
					mu.Unlock()
					return
				}
				mu.Unlock()
				
				// Get next item
				item := queue[0]
				queue = queue[1:] // Dequeue
				
				// Skip if already processed
				if processedElements[item.ElementName] {
					continue
				}
				
				// Get the graph node for this element
				elementNode, exists := GetElementsGraphNodeByName(item.ElementName)
				if !exists || elementNode == nil {
					continue
				}
				
				// If base element or no recipes, create a simple tree
				if IsBaseElement(elementNode.Name) || len(elementNode.RecipesToMakeThisElement) == 0 {
					simpleTree := &RecipeTreeNode{
						Name:      elementNode.Name,
						ImagePath: GetImagePath(elementNode.ImagePath),
					}
					elementToTrees[elementNode.Name] = []*RecipeTreeNode{simpleTree}
					processedElements[elementNode.Name] = true
					continue
				}
				
				// If we don't have trees for all prerequisites, add back to queue with lower priority
				if item.TreeSoFar == nil {
					// Process each recipe for this element
					var allPrereqsProcessed = true
					
					// Check if all recipes are ready to be processed
					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						if !processedElements[elementRecipe.ElementOne.Name] || 
						   !processedElements[elementRecipe.ElementTwo.Name] {
							allPrereqsProcessed = false
							
							// Enqueue prerequisites with higher priority (lower level)
							if !processedElements[elementRecipe.ElementOne.Name] {
								queue = append(queue, &QueueItem{
									ElementName: elementRecipe.ElementOne.Name, 
									Level:       item.Level + 1,
								})
							}
							
							if !processedElements[elementRecipe.ElementTwo.Name] {
								queue = append(queue, &QueueItem{
									ElementName: elementRecipe.ElementTwo.Name, 
									Level:       item.Level + 1,
								})
							}
						}
					}
					
					// If not all prerequisites are processed, re-add this item with lower priority
					if !allPrereqsProcessed {
						queue = append(queue, &QueueItem{
							ElementName: item.ElementName,
							Level:       item.Level + 10, // Lower priority
						})
						continue
					}
					
					// All prerequisites are processed, create trees for this element
					var elementTrees []*RecipeTreeNode
					
					// Create trees for each recipe of this element
					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						leftTrees := elementToTrees[elementRecipe.ElementOne.Name]
						rightTrees := elementToTrees[elementRecipe.ElementTwo.Name]
						
						// Combine all possible combinations
						for _, lt := range leftTrees {
							for _, rt := range rightTrees {
								// Create a new tree with this recipe
								newTree := &RecipeTreeNode{
									Name:      elementNode.Name,
									ImagePath: GetImagePath(elementNode.ImagePath),
									Element1:  lt,
									Element2:  rt,
								}

								// Signal tree change
								if signalTreeChange != nil {
									signalTreeChange(newTree)
								}

								elementTrees = append(elementTrees, newTree)
							}
						}
					}
					
					elementToTrees[elementNode.Name] = elementTrees
					processedElements[elementNode.Name] = true
				}
			}
			
			// After processing all elements, combine the trees for the original recipe
			leftTrees := elementToTrees[r.ElementOne.Name]
			rightTrees := elementToTrees[r.ElementTwo.Name]
			
			if leftTrees == nil || rightTrees == nil {
				return // Missing trees for required elements
			}
			
			// Create all possible combinations
			for _, lt := range leftTrees {
				for _, rt := range rightTrees {
					mu.Lock()
					if treesFound >= maxTreeCount {
						mu.Unlock()
						return
					}
					
					// Create final tree
					root := &RecipeTreeNode{
						Name:      targetGraphNode.Name,
						ImagePath: GetImagePath(targetGraphNode.ImagePath),
						Element1:  lt,
						Element2:  rt,
					}
					
					treesFound++
					mu.Unlock()

					// Signal tree change
					if signalTreeChange != nil {
						signalTreeChange(root)
					}
					
					resultChan <- root
				}
			}
		}(recipe)
	}
	
	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results
	for tree := range resultChan {
		result = append(result, tree)
		if len(result) >= maxTreeCount {
			break
		}
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("no complete tree found for target %s", targetGraphNode.Name)
	}
	
	return result, nil
}