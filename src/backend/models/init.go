package models

func Init() {
	InitElementsGraph()
}

func InitElementsGraph() {
	elements, err := LoadElementsFromJSON("./data/elements.json")
	if err != nil {
		panic(err)
	}

	// Step 1: Create nodes for each element
	for _, el := range elements {
		node := &ElementsGraphNode{
			Name:                      el.Name,
			ImagePath:                 el.ImagePath,
			RecipesToMakeThisElement:  []*Recipe{},
			RecipesToMakeOtherElement: []*Recipe{},
			IsVisited:                 false,
		}
		nameToNode[el.Name] = node
	}

	// Step 2: Populate all RecipesToMakeThisElement and RecipesToMakeOtherElement
	for _, el := range elements {
		resultNode := nameToNode[el.Name]
		for _, r := range el.Recipes {
			if len(r) != 2 {
				continue
			}
			ing1, ing2 := r[0], r[1]
			node1, ok1 := nameToNode[ing1]
			node2, ok2 := nameToNode[ing2]
			if !ok1 || !ok2 {
				continue
			}

			// Check for duplicates before appending
			recipe := &Recipe{
				ElementOne: node1,
				ElementTwo: node2,
			}

			// Only add the recipe if it's not already present
			// Avoid duplicate recipes for resultNode
			if !containsRecipe(resultNode.RecipesToMakeThisElement, recipe) {
				resultNode.RecipesToMakeThisElement = append(resultNode.RecipesToMakeThisElement, recipe)
			}
			if !containsRecipe(node1.RecipesToMakeOtherElement, recipe) {
				node1.RecipesToMakeOtherElement = append(node1.RecipesToMakeOtherElement, recipe)
			}
			if !containsRecipe(node2.RecipesToMakeOtherElement, recipe) {
				node2.RecipesToMakeOtherElement = append(node2.RecipesToMakeOtherElement, recipe)
			}
		}
	}

	// Step 3: Add basic elements to root node
	basics := []string{"Air", "Earth", "Fire", "Water"}
	for _, name := range basics {
		if node, ok := nameToNode[name]; ok {
			ElementsGraph.RecipesToMakeOtherElement = append(ElementsGraph.RecipesToMakeOtherElement, &Recipe{
				ElementOne: node,
				ElementTwo: nil,
			})
		}
	}

	// Step 4: Traverse using DFS to establish connections
	visited := make(map[string]bool)

	var dfs func(node *ElementsGraphNode)
	dfs = func(node *ElementsGraphNode) {
		if visited[node.Name] {
			return
		}
		visited[node.Name] = true
		node.IsVisited = true

		// For each recipe that uses this element
		for _, recipe := range node.RecipesToMakeOtherElement {
			// Skip if recipe is incomplete
			if recipe.ElementTwo == nil {
				continue
			}

			// Check if both ingredients of this recipe are visited
			if !recipe.ElementOne.IsVisited || !recipe.ElementTwo.IsVisited {
				continue
			}

			// Find the resulting element for this recipe
			for _, potential := range nameToNode {
				for _, r := range potential.RecipesToMakeThisElement {
					if recipesMatch(r, recipe) {
						dfs(potential)
					}
				}
			}
		}
	}

	// Start DFS from basic elements
	for _, name := range basics {
		if node, ok := nameToNode[name]; ok {
			dfs(node)
		}
	}

}

func containsRecipe(recipes []*Recipe, recipe *Recipe) bool {
	for _, r := range recipes {
		if r.ElementOne == recipe.ElementOne && r.ElementTwo == recipe.ElementTwo {
			return true
		}
	}
	return false
}
