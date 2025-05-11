package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func BFSFindTrees(
	targetGraphNode *ElementsGraphNode,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode, int, int32),
	globalStartTime time.Time,
	globalNodeCounter *int32,
	delayMs int,
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
		result     []*RecipeTreeNode
		treesFound int
		mu         sync.Mutex
		wg         sync.WaitGroup
		resultChan = make(chan *RecipeTreeNode, maxTreeCount)
	)

	type QueueItem struct {
		ElementName string
		Level       int
		TreeSoFar   *RecipeTreeNode
		IsComplete  bool
	}

	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			elementToTrees := make(map[string][]*RecipeTreeNode)
			processedElements := make(map[string]bool)
			queue := make([]*QueueItem, 0)
			queue = append(
				queue,
				&QueueItem{ElementName: r.ElementOne.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
				&QueueItem{ElementName: r.ElementTwo.Name, Level: 1, TreeSoFar: nil, IsComplete: false},
			)

			for len(queue) > 0 {
				if delayMs > 0 {
					time.Sleep(time.Duration(delayMs) * time.Millisecond)
				}

				mu.Lock()
				if treesFound >= maxTreeCount {
					mu.Unlock()
					return
				}
				mu.Unlock()

				item := queue[0]
				queue = queue[1:]

				if processedElements[item.ElementName] {
					continue
				}

				elementNode, exists := GetElementsGraphNodeByName(item.ElementName)
				if !exists || elementNode == nil {
					continue
				}

				atomic.AddInt32(globalNodeCounter, 1)

				if IsBaseElement(elementNode.Name) || len(elementNode.RecipesToMakeThisElement) == 0 {
					simpleTree := &RecipeTreeNode{
						Name:      elementNode.Name,
						ImagePath: GetImagePath(elementNode.ImagePath),
					}
					elementToTrees[elementNode.Name] = []*RecipeTreeNode{simpleTree}
					processedElements[elementNode.Name] = true
					continue
				}

				if item.TreeSoFar == nil {
					var allPrereqsProcessed = true

					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						if !processedElements[elementRecipe.ElementOne.Name] ||
							!processedElements[elementRecipe.ElementTwo.Name] {
							allPrereqsProcessed = false

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

					if !allPrereqsProcessed {
						queue = append(queue, &QueueItem{
							ElementName: item.ElementName,
							Level:       item.Level + 10,
						})
						continue
					}

					var elementTrees []*RecipeTreeNode
					for _, elementRecipe := range elementNode.RecipesToMakeThisElement {
						leftTrees := elementToTrees[elementRecipe.ElementOne.Name]
						rightTrees := elementToTrees[elementRecipe.ElementTwo.Name]

						for _, lt := range leftTrees {
							for _, rt := range rightTrees {
								newTree := &RecipeTreeNode{
									Name:      elementNode.Name,
									ImagePath: GetImagePath(elementNode.ImagePath),
									Element1:  lt,
									Element2:  rt,
								}

								if signalTreeChange != nil {
									func() {
										defer func() {
											if r := recover(); r != nil {
											}
										}()
										signalTreeChange(
											newTree,
											int(time.Since(globalStartTime).Milliseconds()),
											atomic.LoadInt32(globalNodeCounter),
										)
									}()
								}

								elementTrees = append(elementTrees, newTree)
							}
						}
					}

					elementToTrees[elementNode.Name] = elementTrees
					processedElements[elementNode.Name] = true
				}
			}

			leftTrees := elementToTrees[r.ElementOne.Name]
			rightTrees := elementToTrees[r.ElementTwo.Name]

			if leftTrees == nil || rightTrees == nil {
				return
			}

			for _, lt := range leftTrees {
				for _, rt := range rightTrees {
					mu.Lock()
					if treesFound >= maxTreeCount {
						mu.Unlock()
						return
					}

					root := &RecipeTreeNode{
						Name:      targetGraphNode.Name,
						ImagePath: GetImagePath(targetGraphNode.ImagePath),
						Element1:  lt,
						Element2:  rt,
					}

					treesFound++
					mu.Unlock()

					if signalTreeChange != nil {
						func() {
							defer func() {
								if r := recover(); r != nil {
								}
							}()
							signalTreeChange(
								root,
								int(time.Since(globalStartTime).Milliseconds()),
								atomic.LoadInt32(globalNodeCounter),
							)
						}()
					}

					resultChan <- root
				}
			}
		}(recipe)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

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
