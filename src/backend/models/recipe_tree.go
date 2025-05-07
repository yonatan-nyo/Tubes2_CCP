package models

import "errors"

// Struct RecipeTreeNode
type RecipeTreeNode struct {
	Name      string          `json:"name"`
	ImagePath string          `json:"image_path"`
	Element1  *RecipeTreeNode `json:"element_1,omitempty"`
	Element2  *RecipeTreeNode `json:"element_2,omitempty"`
}

// Check if the element is a base element
// Base elements: Air, Water, Earth, Fire
func isBaseElement(name string) bool {
	return name == "Air" || name == "Water" || name == "Earth" || name == "Fire"
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


func GetRecipeTreeBidirectional(target string) (*RecipeTreeNode, error) {
	targetNode, ok := nameToNode[target]
	if !ok {
		return nil, ErrElementNotFound(target)
	}

	if isBaseElement(target) {
		return &RecipeTreeNode{
			Name:      target,
			ImagePath: GetImagePath(targetNode.ImagePath),
		}, nil
	}

	parentRecipes := make(map[string]*Recipe)

	forwardQueue := []*ElementsGraphNode{targetNode}
	forwardVisited := map[string]bool{target: true}

	backwardVisited := map[string]bool{
		"Air": true, "Water": true, "Earth": true, "Fire": true,
	}

	distanceFromBase := map[string]int{
		"Air": 0, "Water": 0, "Earth": 0, "Fire": 0,
	}

	baseQueue := []*ElementsGraphNode{}
	for _, name := range []string{"Air", "Water", "Earth", "Fire"} {
		if baseNode, exists := nameToNode[name]; exists {
			baseQueue = append(baseQueue, baseNode)
		}
	}

	for len(baseQueue) > 0 {
		node := baseQueue[0]
		baseQueue = baseQueue[1:]

		currentDistance := distanceFromBase[node.Name]
		nextDistance := currentDistance + 1

		for _, recipe := range node.RecipesToMakeOtherElement {
			resultName := recipe.TargetElementName
			resultNode, exists := nameToNode[resultName]
			if !exists {
				continue
			}

			if recipe.ElementOne != nil && (recipe.ElementTwo != nil || isBaseElement(recipe.ElementOne.Name)) {
				if _, visited := distanceFromBase[resultName]; !visited || nextDistance < distanceFromBase[resultName] {
					distanceFromBase[resultName] = nextDistance
					backwardVisited[resultName] = true
					parentRecipes[resultName] = recipe
					baseQueue = append(baseQueue, resultNode)
				}
			}
		}
	}

	var connectionNode string
	minDistance := -1

	for len(forwardQueue) > 0 {
		node := forwardQueue[0]
		forwardQueue = forwardQueue[1:]

		if backwardVisited[node.Name] {
			forwardDistance := 0 // Not tracked
			backwardDistance := distanceFromBase[node.Name]
			totalDistance := forwardDistance + backwardDistance

			if minDistance == -1 || totalDistance < minDistance {
				minDistance = totalDistance
				connectionNode = node.Name
			}
		}

		for _, recipe := range node.RecipesToMakeThisElement {
			if recipe.ElementOne == nil || (recipe.ElementTwo == nil && !isBaseElement(recipe.ElementOne.Name)) {
				continue
			}

			for _, ingredient := range []*ElementsGraphNode{recipe.ElementOne, recipe.ElementTwo} {
				if ingredient == nil {
					continue
				}

				ingredientName := ingredient.Name
				if !forwardVisited[ingredientName] {
					forwardVisited[ingredientName] = true
					forwardQueue = append(forwardQueue, ingredient)

					if _, exists := parentRecipes[ingredientName]; !exists {
						var otherElement *ElementsGraphNode
						if recipe.ElementOne == ingredient {
							otherElement = recipe.ElementTwo
						} else {
							otherElement = recipe.ElementOne
						}

						parentRecipes[ingredientName] = &Recipe{
							ElementOne:        ingredient,
							ElementTwo:        otherElement,
							TargetElementName: node.Name,
						}
					}
				}
			}
		}
	}

	if connectionNode == "" {
		return nil, errors.New("no path found from base elements to target")
	}

	return buildTreeFromParents(connectionNode, target, parentRecipes), nil
}

func buildTreeFromParents(start, target string, parentRecipes map[string]*Recipe) *RecipeTreeNode {
	treeNodes := make(map[string]*RecipeTreeNode)

	var buildNode func(name string) *RecipeTreeNode
	buildNode = func(name string) *RecipeTreeNode {
		if node, exists := treeNodes[name]; exists {
			return node
		}

		elemNode, exists := nameToNode[name]
		if !exists {
			return nil
		}

		treeNode := &RecipeTreeNode{
			Name:      name,
			ImagePath: GetImagePath(elemNode.ImagePath),
		}
		treeNodes[name] = treeNode

		if isBaseElement(name) || name == start {
			return treeNode
		}

		recipe, exists := parentRecipes[name]
		if !exists {
			return treeNode
		}

		if recipe.ElementOne != nil {
			treeNode.Element1 = buildNode(recipe.ElementOne.Name)
		}
		if recipe.ElementTwo != nil {
			treeNode.Element2 = buildNode(recipe.ElementTwo.Name)
		}

		return treeNode
	}

	return buildNode(target)
}
