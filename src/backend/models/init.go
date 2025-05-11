package models

import (
	"fmt"
	"slices"
)

func Init() {
	InitElementsGraph()
}

var baseElements = []string{"Air", "Earth", "Fire", "Water"}

func InitElementsGraph() {
	elements, err := LoadElementsFromJSON("./data/elements.json")
	if err != nil {
		panic(err)
	}

	// Initialize the left side of the table (target-recipe) the target element
	for _, el := range elements {
		node := &ElementsGraphNode{
			Name:                      el.Name,
			ImagePath:                 el.ImagePath,
			RecipesToMakeThisElement:  []*Recipe{},
			RecipesToMakeOtherElement: []*Recipe{},
			IsVisited:                 false,
			Tier:                      -1,
		}
		nameToNode[el.Name] = node
	}

	var elementsNameNotFound []string
	for _, el := range elements {
		if _, ok := nameToNode[el.Name]; !ok {
			elementsNameNotFound = append(elementsNameNotFound, el.Name)
		}
		for _, r := range el.Recipes {
			if len(r) != 2 {
				continue
			}
			ing1, ing2 := r[0], r[1]
			_, ok1 := nameToNode[ing1]
			if !ok1 {
				elementsNameNotFound = append(elementsNameNotFound, ing1)
			}
			_, ok2 := nameToNode[ing2]
			if !ok2 {
				elementsNameNotFound = append(elementsNameNotFound, ing2)
			}
		}
	}
	
	for _, node := range nameToNode {
		node.MadeFrom = make(map[string]bool)
	}

	// Output the elements that are not found
	if len(elementsNameNotFound) > 0 {
		fmt.Println("Elements not found in the graph:")
		for _, el := range elementsNameNotFound {
			fmt.Println(el)
		}
	} else {
		fmt.Println("All elements found in the graph")
	}

	// Populate all RecipesToMakeThisElement and RecipesToMakeOtherElement
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

	// Add basic elements to root node
	basics := []string{"Air", "Earth", "Fire", "Water"}
	for _, name := range basics {
		if node, ok := nameToNode[name]; ok {
			ElementsGraph.RecipesToMakeOtherElement = append(ElementsGraph.RecipesToMakeOtherElement, &Recipe{
				ElementOne: node,
				ElementTwo: nil,
			})
		}
	}

	// For every element that doesnt have recipe to make this element, append to basics
	for _, node := range nameToNode {
		if node.Name == "Air" || node.Name == "Earth" || node.Name == "Fire" || node.Name == "Water" {
			node.Tier = 0
			continue
		}
		if len(node.RecipesToMakeThisElement) == 0 {
			ElementsGraph.RecipesToMakeOtherElement = append(ElementsGraph.RecipesToMakeOtherElement, &Recipe{
				ElementOne: node,
				ElementTwo: nil,
			})
			node.Tier = 0
			baseElements = append(baseElements, node.Name)
		}
	}

	//  the tier for each node
	curTier := 0

	for {
		nodesWithInitializedTier := map[string]bool{}
		for _, node := range nameToNode {
			if node.Tier != -1 {
				nodesWithInitializedTier[node.Name] = true
			}
		}

		if len(nodesWithInitializedTier) == len(nameToNode) {
			break
		}

		// Set the tier for the next level
		for _, node := range nameToNode {
			if node.Tier != -1 {
				continue
			}
			//check if the recipe to make this element is feasible
			for _, recipe := range node.RecipesToMakeThisElement {
				if nodesWithInitializedTier[recipe.ElementOne.Name] && (recipe.ElementTwo == nil || nodesWithInitializedTier[recipe.ElementTwo.Name]) {
					node.Tier = curTier + 1
					break
				}
			}
		}
		curTier++
	}

	// Populate MadeFrom
	for _, node := range nameToNode {
		for _, recipe := range node.RecipesToMakeThisElement {
			if recipe.ElementOne != nil {
				node.MadeFrom[recipe.ElementOne.Name] = true
				if recipe.ElementOne.MadeFrom != nil {
					for made := range recipe.ElementOne.MadeFrom {
						node.MadeFrom[made] = true
					}
				}
			}
			if recipe.ElementTwo != nil {
				node.MadeFrom[recipe.ElementTwo.Name] = true
				if recipe.ElementTwo.MadeFrom != nil {
					for made := range recipe.ElementTwo.MadeFrom {
						node.MadeFrom[made] = true
					}
				}
			}
		}
	}

	// Element used in a recipe must have lower tier than the target node
	for _, node := range nameToNode {
		// Create a new slice for recipes to keep
		filtered := make([]*Recipe, 0, len(node.RecipesToMakeThisElement))

		// Add only valid recipes to the filtered slice
		for _, recipe := range node.RecipesToMakeThisElement {
			shouldKeep := true

			// Check if any element in the recipe has a tier >= the target node tier
			if recipe.ElementOne != nil && recipe.ElementOne.Tier >= node.Tier {
				shouldKeep = false
			}
			if recipe.ElementTwo != nil && recipe.ElementTwo.Tier >= node.Tier {
				shouldKeep = false
			}

			// Keep only valid recipes
			if shouldKeep {
				filtered = append(filtered, recipe)
			}
		}

		// Use slices to create a new slice (just to keep the import)
		node.RecipesToMakeThisElement = slices.Clone(filtered)
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

func IsBaseElement(name string) bool {
	for _, base := range baseElements {
		if name == base {
			return true
		}
	}
	return false
}
