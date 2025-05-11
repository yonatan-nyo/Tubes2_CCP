package models

import (
	"fmt"
	"time"
)

type RecipeTreeNode struct {
	Name      string          `json:"name"`
	ImagePath string          `json:"image_path"`
	Element1  *RecipeTreeNode `json:"element_1,omitempty"`
	Element2  *RecipeTreeNode `json:"element_2,omitempty"`
}

func ValidateInputParams(
	target string,
	mode string,
	maxTreeCount int,
) error {
	if maxTreeCount <= 0 {
		return fmt.Errorf("maxTreeCount must be greater than 0")
	}

	if nameToNode == nil {
		return fmt.Errorf("elements graph is not initialized")
	}

	targetGraphNode, ok := nameToNode[target]
	if !ok || targetGraphNode == nil {
		return fmt.Errorf("target %s not found or is nil in elements graph", target)
	}

	return nil
}

func GenerateRecipeTree(
	target string,
	mode string,
	maxTreeCount int,
	signallerFn func(*RecipeTreeNode, int, int32),
	delayMs int,
	globalStartTime time.Time,
	globalNodeCount *int32,
) ([]*RecipeTreeNode, error) {
	if err := ValidateInputParams(target, mode, maxTreeCount); err != nil {
		return nil, err
	}

	rootRecipeTree := &RecipeTreeNode{
		Name:      target,
		ImagePath: GetImagePath(target),
	}

	targetGraphNode, ok := nameToNode[target]
	if !ok || targetGraphNode == nil {
		return nil, fmt.Errorf("target %s not found or is nil in elements graph", target)
	}

	var (
		trees []*RecipeTreeNode
		err   error
	)

	if trees, err = ProcessRecipeTree(
		rootRecipeTree,
		targetGraphNode,
		mode,
		maxTreeCount,
		signallerFn,
		globalStartTime,
		delayMs,
		globalNodeCount,
	); err != nil {
		return nil, err
	}

	if len(trees) == 0 {
		return nil, fmt.Errorf("no complete tree found for target %s", target)
	}

	return trees, nil
}

func ProcessRecipeTree(
	rootRecipeTree *RecipeTreeNode,
	targetGraphNode *ElementsGraphNode,
	mode string,
	maxTreeCount int,
	signalTreeChange func(*RecipeTreeNode, int, int32),
	globalStartTime time.Time,
	delayMs int,
	globalNodeCounter *int32,
) ([]*RecipeTreeNode, error) {

	if mode == "dfs" {
		return DFSFindTrees(
			rootRecipeTree,
			targetGraphNode,
			maxTreeCount,
			signalTreeChange,
			globalStartTime,
			globalNodeCounter,
			delayMs,
		)
	}
	if mode == "bfs" {
		return BFSFindTrees(
			targetGraphNode,
			maxTreeCount,
			signalTreeChange,
			globalStartTime,
			globalNodeCounter,
			delayMs,
		)
	}
	if mode == "bidirectional" {
		return BidirectionalFindTrees(
			rootRecipeTree,
			targetGraphNode,
			maxTreeCount,
			signalTreeChange,
			globalStartTime,
			globalNodeCounter,
			delayMs,
		)
	}

	return nil, fmt.Errorf("invalid mode: %s", mode)
}
