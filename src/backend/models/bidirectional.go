package models

import (
	"fmt"
	"time"
	//"sync/atomic"
)

type QueueItem struct {
	Element *ElementsGraphNode
	From    string
}

func BidirectionalFindTrees(
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

	visitedUpper := make(map[string]bool)
	visitedLower := make(map[string]bool)
	seenMeeting := make(map[string]bool)

	queueUpper := []*QueueItem{{Element: targetGraphNode}}
	queueLower := []*QueueItem{}
	for _, base := range GetBaseElements() {
		if targetGraphNode.IsThisMadeFrom(base) {
			if node, ok := GetElementsGraphNodeByName(base); ok {
				queueLower = append(queueLower, &QueueItem{Element: node})
			}
		}
	}

	var resultTrees []*RecipeTreeNode

	for len(queueUpper) > 0 && len(queueLower) > 0 && len(resultTrees) < maxTreeCount {
		nQueueUpper, newUpperNames := processUpper(queueUpper, visitedUpper)
		for _, name := range newUpperNames {
			if visitedLower[name] && name == targetGraphNode.Name {
				if _, already := seenMeeting[name]; already {
					continue
				}
				seenMeeting[name] = true
				if node, ok := GetElementsGraphNodeByName(name); ok {
					treesFromDFS, err := DFSFindTrees(nil, node, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
					if err == nil {
						resultTrees = appendAllValidTargetTrees(resultTrees, treesFromDFS, targetGraphNode.Name, maxTreeCount)
						if len(resultTrees) >= maxTreeCount {
							return resultTrees, nil
						}
					}
				}
			}
		}
		queueUpper = nQueueUpper

		nQueueLower, newLowerNames := processLower(queueLower, visitedLower)
		for _, name := range newLowerNames {
			if visitedUpper[name] && name == targetGraphNode.Name {
				if _, already := seenMeeting[name]; already {
					continue
				}
				seenMeeting[name] = true
				if node, ok := GetElementsGraphNodeByName(name); ok {
					treesFromDFS, err := DFSFindTrees(nil, node, maxTreeCount, signalTreeChange, globalStartTime, globalNodeCounter, delayMs)
					if err == nil {
						resultTrees = appendAllValidTargetTrees(resultTrees, treesFromDFS, targetGraphNode.Name, maxTreeCount)
						if len(resultTrees) >= maxTreeCount {
							return resultTrees, nil
						}
					}
				}
			}
		}
		queueLower = nQueueLower
	}
	return resultTrees, nil
}

func processUpper(queue []*QueueItem, visited map[string]bool) ([]*QueueItem, []string) {
	nextQueue := []*QueueItem{}
	produced := []string{}
	for _, item := range queue {
		node := item.Element
		if visited[node.Name] {
			continue
		}
		visited[node.Name] = true
		produced = append(produced, node.Name)
		for _, recipe := range node.RecipesToMakeThisElement {
			nextQueue = append(nextQueue, &QueueItem{Element: recipe.ElementOne})
			nextQueue = append(nextQueue, &QueueItem{Element: recipe.ElementTwo})
		}
	}
	return nextQueue, produced
}

func processLower(queue []*QueueItem, visited map[string]bool) ([]*QueueItem, []string) {
	nextQueue := []*QueueItem{}
	produced := []string{}
	for _, item := range queue {
		node := item.Element
		if visited[node.Name] {
			continue
		}
		visited[node.Name] = true
		produced = append(produced, node.Name)
		for _, recipe := range node.RecipesToMakeOtherElement {
			if targetNode, ok := GetElementsGraphNodeByName(recipe.TargetElementName); ok {
				nextQueue = append(nextQueue, &QueueItem{Element: targetNode})
			}
		}
	}
	return nextQueue, produced
}

func appendAllValidTargetTrees(existing []*RecipeTreeNode, newTrees []*RecipeTreeNode, target string, maxTreeCount int) []*RecipeTreeNode {
	existingSet := make(map[string]bool)
	for _, t := range existing {
		if t.Name == target {
			existingSet[t.ImagePath+t.Name+treeToString(t)] = true
		}
	}
	for _, t := range newTrees {
		if t.Name == target {
			key := t.ImagePath + t.Name + treeToString(t)
			if !existingSet[key] {
				existing = append(existing, t)
				existingSet[key] = true
				if len(existing) >= maxTreeCount {
					break
				}
			}
		}
	}
	return existing
}

func treeToString(tree *RecipeTreeNode) string {
	if tree == nil {
		return ""
	}
	return tree.Name + "(" + treeToString(tree.Element1) + "," + treeToString(tree.Element2) + ")"
}


// Helper function to return base elements
func GetBaseElements() []string {
    return baseElements
}