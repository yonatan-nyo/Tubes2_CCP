package models

import (
	"fmt"
	"sync"
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
		wg     sync.WaitGroup
		count  = 0
	)

	// Buffered channel to limit creation to maxTreeCount
	treeChan := make(chan *RecipeTreeNode, maxTreeCount)

	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			leftTrees, err1 := DFSFindTrees(nil, r.ElementOne, maxTreeCount, signalTreeChange)
			if err1 != nil {
				return
			}

			rightTrees, err2 := DFSFindTrees(nil, r.ElementTwo, maxTreeCount, signalTreeChange)
			if err2 != nil {
				return
			}

			for _, lt := range leftTrees {
				mu.Lock()
				if count >= maxTreeCount {
					mu.Unlock()
					break
				}
				mu.Unlock()

				for _, rt := range rightTrees {
					mu.Lock()
					if count >= maxTreeCount {
						mu.Unlock()
						break
					}
					root := &RecipeTreeNode{
						Name:      targetGraphNode.Name,
						ImagePath: GetImagePath(targetGraphNode.ImagePath),
						Element1:  lt,
						Element2:  rt,
					}
					treeChan <- root
					count++
					mu.Unlock()
				}
			}
		}(recipe)

		mu.Lock()
		if count >= maxTreeCount {
			mu.Unlock()
			break
		}
		mu.Unlock()
	}

	go func() {
		wg.Wait()
		close(treeChan)
	}()

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
