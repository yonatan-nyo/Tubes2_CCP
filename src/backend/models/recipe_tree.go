package models

import "errors"

// Struct RecipeTreeNode
type RecipeTreeNode struct {
	Name      string          `json:"name"`
	ImagePath string          `json:"image_path"`
	Element1  *RecipeTreeNode `json:"element_1,omitempty"`
	Element2  *RecipeTreeNode `json:"element_2,omitempty"`
}

// ErrElementNotFound: element not found in graph
func ErrElementNotFound(name string) error {
	return errors.New("element not found: " + name)
}

func GetRecipeTreeDFS(target string) (*RecipeTreeNode, error) {
	targetNode, exists := nameToNode[target]
	if !exists {
		return nil, ErrElementNotFound(target)
	}
	visited := make(map[string]bool)
	return buildRecipeTreeDFS(targetNode, visited), nil
}

func buildRecipeTreeDFS(node *ElementsGraphNode, visited map[string]bool) *RecipeTreeNode {
	if node == nil {
		return nil
	}

	if isBaseElement(node.Name) {
		return &RecipeTreeNode{Name: node.Name, ImagePath: GetImagePath(node.ImagePath)}
	}

	if visited[node.Name] {
		return nil
	}
	visited[node.Name] = true

	if len(node.RecipesToMakeThisElement) == 0 {
		return &RecipeTreeNode{Name: node.Name, ImagePath: GetImagePath(node.ImagePath)}
	}

	best := node.RecipesToMakeThisElement[0]

	tree := &RecipeTreeNode{
		Name:      node.Name,
		ImagePath: GetImagePath(node.ImagePath),
		Element1:  buildRecipeTreeDFS(best.ElementOne, visited),
		Element2:  buildRecipeTreeDFS(best.ElementTwo, visited),
	}

	return tree
}

func GetRecipeTreeBFS(target string) (*RecipeTreeNode, error) {
	rootNode, exists := nameToNode[target]
	if !exists {
		return nil, ErrElementNotFound(target)
	}

	type QueueItem struct {
		Node    *ElementsGraphNode
		TreeRef **RecipeTreeNode
	}
	visited := make(map[string]bool)

	root := &RecipeTreeNode{
		Name:      rootNode.Name,
		ImagePath: GetImagePath(rootNode.ImagePath),
	}
	queue := []QueueItem{{Node: rootNode, TreeRef: &root}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		node := current.Node
		treePtr := current.TreeRef

		if isBaseElement(node.Name) || len(node.RecipesToMakeThisElement) == 0 {
			continue
		}

		best := node.RecipesToMakeThisElement[0]

		if best.ElementOne != nil {
			elem1 := &RecipeTreeNode{
				Name:      best.ElementOne.Name,
				ImagePath: GetImagePath(best.ElementOne.ImagePath),
			}
			(*treePtr).Element1 = elem1

			if !isBaseElement(elem1.Name) && !visited[elem1.Name] {
				visited[elem1.Name] = true
				queue = append(queue, QueueItem{Node: best.ElementOne, TreeRef: &(*treePtr).Element1})
			}
		}

		if best.ElementTwo != nil {
			elem2 := &RecipeTreeNode{
				Name:      best.ElementTwo.Name,
				ImagePath: GetImagePath(best.ElementTwo.ImagePath),
			}
			(*treePtr).Element2 = elem2

			if !isBaseElement(elem2.Name) && !visited[elem2.Name] {
				visited[elem2.Name] = true
				queue = append(queue, QueueItem{Node: best.ElementTwo, TreeRef: &(*treePtr).Element2})
			}
		}
	}

	return root, nil
}

// Check if the element is a base element
// Base elements: Air, Water, Earth, Fire
func isBaseElement(name string) bool {
	return name == "Air" || name == "Water" || name == "Earth" || name == "Fire"
}