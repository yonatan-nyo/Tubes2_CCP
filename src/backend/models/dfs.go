package models

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func DFSFindTrees(
	rootRecipeTree *RecipeTreeNode,
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
		result []*RecipeTreeNode
		mu     sync.Mutex
		wg     sync.WaitGroup
		count  = 0
	)

	treeChan := make(chan *RecipeTreeNode, maxTreeCount)

	for _, recipe := range targetGraphNode.RecipesToMakeThisElement {
		wg.Add(1)
		go func(r *Recipe) {
			defer wg.Done()

			if delayMs > 0 {
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}

			atomic.AddInt32(globalNodeCounter, 1)

			leftTrees, err1 := DFSFindTrees(nil, r.ElementOne, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
			if err1 != nil {
				return
			}

			rightTrees, err2 := DFSFindTrees(nil, r.ElementTwo, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
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
