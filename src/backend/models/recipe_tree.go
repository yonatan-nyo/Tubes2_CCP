package models

import "fmt"

// Struct RecipeTreeNode
type RecipeTreeNode struct {
	Name                   string          `json:"name"`
	ImagePath              string          `json:"image_path"`
	Element1               *RecipeTreeNode `json:"element_1,omitempty"`
	Element2               *RecipeTreeNode `json:"element_2,omitempty"`
	IsParentElement        map[string]bool `json:"is_parent_element,omitempty"`
	MinimumNodesRecipeTree int             `json:"minimum_nodes_recipe_tree"`
}

// ValidateInputParams validates the input parameters for GetRecipeTree
func ValidateInputParams(
	target string,
	mode string,
	findBestTree bool,
	maxTreeCount int,
) error {
	// if find best tree is true, maxTreeCount should be 0
	if findBestTree && maxTreeCount != 0 {
		return fmt.Errorf("maxTreeCount should be 0 when findBest is true")
	}
	// if find best tree is false, maxTreeCount should be greater than 0
	if !findBestTree && maxTreeCount <= 0 {
		return fmt.Errorf("maxTreeCount should be greater than 0 when findBest is false")
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
	findBestTree bool,
	maxTreeCount int,
	signallerFn func(*RecipeTreeNode, *RecipeTreeNode),
) (*RecipeTreeNode, error) {
	// Validate input parameters
	if err := ValidateInputParams(target, mode, findBestTree, maxTreeCount); err != nil {
		return nil, err
	}

	rootRecipeTree := &RecipeTreeNode{
		Name:      target,
		ImagePath: GetImagePath(target),
	}

	if IsBaseElement(target) {
		rootRecipeTree.MinimumNodesRecipeTree = 1
		return rootRecipeTree, nil
	}

	// process it
	if err := ProcessRecipeTree(
		rootRecipeTree,
		target,
		mode,
		findBestTree,
		maxTreeCount,
		signallerFn,
	); err != nil {
		return nil, err
	}

	return rootRecipeTree, nil
}

func ProcessRecipeTree(
	rootRecipeTree *RecipeTreeNode, // gets passed on both
	target string,
	mode string,
	findBestTree bool,
	maxTreeCount int, //get passed when findBestTree is false
	signalTreeChange func(*RecipeTreeNode, *RecipeTreeNode), // gets passed on both
	// signalTreeChange to signal the tree change to the queue for tree process animation
) error {
	if mode == "dfs" {
		if findBestTree {

			bestTree, err := GenerateDFSFindBestTree(
				target,
				signalTreeChange,
			)
			if err != nil {
				return fmt.Errorf("DFSFindBestTree error: %v", err)
			}
			*rootRecipeTree = *bestTree
			return nil
		} else {
			fmt.Println("DFSFindTreeWithMaxCount not implemented")
			return fmt.Errorf("DFSFindTreeWithMaxCount not implemented")
		}
	}
	if mode == "bfs" {
		if findBestTree {

			bestTree, err := GenerateBFSFindBestTree(
				target,
				signalTreeChange,
			)
			if err != nil {
				return fmt.Errorf("BFSFindBestTree error: %v", err)
			}
			*rootRecipeTree = *bestTree
			return nil
		} else {
			fmt.Println("BFSFindTreeWithMaxCount not implemented")
			return fmt.Errorf("BFSFindTreeWithMaxCount not implemented")
		}
	}
	if mode == "bidirectional" {
		if findBestTree {
			fmt.Println("BidirectionalFindBestTree not implemented")
			return fmt.Errorf("BidirectionalFindBestTree not implemented")
		} else {
			fmt.Println("BidirectionalFindTreeWithMaxCount not implemented")
			return fmt.Errorf("BidirectionalFindTreeWithMaxCount not implemented")
		}
	}

	return nil
}

func isTreeComplete(
	solutionRecipeTreeNode *RecipeTreeNode,
) bool {
	// Make sure if the leaf node is a base element
	if solutionRecipeTreeNode.Element1 == nil && solutionRecipeTreeNode.Element2 == nil {
		return IsBaseElement(solutionRecipeTreeNode.Name)
	}
	// Check if both elements are base elements
	if solutionRecipeTreeNode.Element1 != nil && solutionRecipeTreeNode.Element2 != nil {
		return IsBaseElement(solutionRecipeTreeNode.Element1.Name) && IsBaseElement(solutionRecipeTreeNode.Element2.Name)
	}
	// Check if one of the elements is a base element
	if solutionRecipeTreeNode.Element1 != nil {
		return IsBaseElement(solutionRecipeTreeNode.Element1.Name)
	}
	if solutionRecipeTreeNode.Element2 != nil {
		return IsBaseElement(solutionRecipeTreeNode.Element2.Name)
	}
	return false
}

/* Clone the tree to avoid modifying the original during traversal
   This is a deep copy function to ensure the original tree remains unchanged */
func cloneTree(node *RecipeTreeNode) *RecipeTreeNode {
	if node == nil {
		return nil
	}
	return &RecipeTreeNode{
		Name:                   node.Name,
		ImagePath:              node.ImagePath,
		Element1:               cloneTree(node.Element1),
		Element2:               cloneTree(node.Element2),
		IsParentElement:        node.IsParentElement,
		MinimumNodesRecipeTree: node.MinimumNodesRecipeTree,
	}
}

// if the recipe element is in isParentElement map of the node
// its gonna
func isMakingCycle(node *RecipeTreeNode, recipe *Recipe) bool {
	if node.IsParentElement == nil {
		node.IsParentElement = make(map[string]bool)
	}

	// check if the recipe is in the isParentElement map
	if _, ok := node.IsParentElement[recipe.TargetElementName]; ok {
		return true
	}

	return false
}
