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

	// Step 2: Find all implicit base elements
	// Elements those are used as ingredients but have no recipes to make them
	ingredientSet := make(map[string]bool) 
	resultSet := make(map[string]bool)    

	for _, el := range elements {
		resultSet[el.Name] = true
		for _, r := range el.Recipes {
			if len(r) != 2 {
				continue
			}
			ingredientSet[r[0]] = true
			ingredientSet[r[1]] = true
		}
	}

	for ing := range ingredientSet {
		if _, exists := nameToNode[ing]; !exists {
			continue
		}
		node := nameToNode[ing]
		if len(node.RecipesToMakeThisElement) == 0 {
			ElementsGraph.RecipesToMakeOtherElement = append(ElementsGraph.RecipesToMakeOtherElement, &Recipe{
				ElementOne: node,
				ElementTwo: nil,
			})
		}
	}

	// Step 3: Populate all RecipesToMakeThisElement and RecipesToMakeOtherElement
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
				ElementOne:        node1,
				ElementTwo:        node2,
				TargetElementName: resultNode.Name,
			}

			// Only add the recipe if it's not already present
			// Avoid duplicate recipes for resultNode
			if !containsRecipe(resultNode.RecipesToMakeThisElement, recipe) {
				resultNode.RecipesToMakeThisElement = append(resultNode.RecipesToMakeThisElement, recipe)
			}
			if !containsRecipe(node1.RecipesToMakeOtherElement, recipe) {
				node1.RecipesToMakeOtherElement = append(node1.RecipesToMakeOtherElement, recipe)
			}
			// Hanya tambahkan ke node2 jika berbeda
			if node1 != node2 && !containsRecipe(node2.RecipesToMakeOtherElement, recipe) {
				node2.RecipesToMakeOtherElement = append(node2.RecipesToMakeOtherElement, recipe)
			}
		}
	}

	// Step 4: Add basic elements to root node
	basics := []string{"Air", "Earth", "Fire", "Water"}
	for _, name := range basics {
		if node, ok := nameToNode[name]; ok {
			ElementsGraph.RecipesToMakeOtherElement = append(ElementsGraph.RecipesToMakeOtherElement, &Recipe{
				ElementOne: node,
				ElementTwo: nil,
			})
		}
	}
}

func containsRecipe(recipes []*Recipe, recipe *Recipe) bool {
	for _, r := range recipes {
		// Handle nil cases properly
		// Case 1: Both recipes have ElementTwo set
		if r.ElementTwo != nil && recipe.ElementTwo != nil {
			sameDirect := r.ElementOne.Name == recipe.ElementOne.Name && r.ElementTwo.Name == recipe.ElementTwo.Name
			sameReverse := r.ElementOne.Name == recipe.ElementTwo.Name && r.ElementTwo.Name == recipe.ElementOne.Name
			if (sameDirect || sameReverse) && r.TargetElementName == recipe.TargetElementName {
				return true
			}
		} else if r.ElementTwo == nil && recipe.ElementTwo == nil {
			// Case 2: Both recipes have ElementTwo as nil
			if r.ElementOne.Name == recipe.ElementOne.Name && r.TargetElementName == recipe.TargetElementName {
				return true
			}
		} else {
			// Case 3: One has ElementTwo as nil and the other doesn't
			// These are different recipes
			continue
		}
	}
	return false
}