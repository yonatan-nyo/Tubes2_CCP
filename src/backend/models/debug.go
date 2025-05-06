package models

import (
	"fmt"
	"strings"
)

func Debug(root *ElementsGraphNode) {
	DebugBasicElementsFromRoot()
	fmt.Println("\n=== Elements Graph Debug Output ===")
	visited := make(map[string]bool)
	for _, recipe := range ElementsGraph.RecipesToMakeOtherElement {
		printNode(recipe.ElementOne, visited, 0)
	}
}

func DebugBasicElementsFromRoot() {
	fmt.Println("=== Basic Elements from Root Node ===")
	for _, recipe := range ElementsGraph.RecipesToMakeOtherElement {
		if recipe.ElementTwo == nil {
			fmt.Printf("- %s (%s)\n", recipe.ElementOne.Name, recipe.ElementOne.ImagePath)
		}
	}
}

// Helper function to recursively print the node and its relationships
// Helper function to recursively print the node and its relationships
func printNode(node *ElementsGraphNode, visited map[string]bool, depth int) {
	if node == nil || visited[node.Name] {
		return
	}
	visited[node.Name] = true

	// Indentation for better readability
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s- Element: %s (%s)\n", indent, node.Name, node.ImagePath)

	// Print the recipes to make this element
	if len(node.RecipesToMakeThisElement) > 0 {
		fmt.Printf("%s  Recipes to make this:\n", indent)
		for _, r := range node.RecipesToMakeThisElement {
			if r.ElementTwo != nil {
				fmt.Printf("%s    %s + %s => %s\n", indent, r.ElementOne.Name, r.ElementTwo.Name, node.Name)
			} else {
				fmt.Printf("%s    %s => %s\n", indent, r.ElementOne.Name, node.Name)
			}
		}
	}

	// Print the recipes using this element to make others
	if len(node.RecipesToMakeOtherElement) > 0 {
		fmt.Printf("%s  Recipes using this element to make others:\n", indent)
		for _, r := range node.RecipesToMakeOtherElement {
			// Directly print the result of the recipe
			targetNode := r.ElementOne
			if r.ElementTwo != nil {
				targetNode = r.ElementTwo
			}
			if targetNode != nil {
				fmt.Printf("%s    %s + %s => %s\n", indent, node.Name, targetNode.Name, targetNode.Name)
			}
		}
	}

	// Recursively print connected nodes based on the recipes
	for _, r := range node.RecipesToMakeOtherElement {
		var target *ElementsGraphNode
		if r.ElementTwo != nil {
			target = r.ElementTwo
		} else {
			target = r.ElementOne
		}
		if target != nil {
			printNode(target, visited, depth+1)
		}
	}
}

// Helper to compare recipes
func recipesMatch(a, b *Recipe) bool {
	if a == nil || b == nil {
		return false
	}
	return a.ElementOne.Name == b.ElementOne.Name &&
		((a.ElementTwo == nil && b.ElementTwo == nil) ||
			(a.ElementTwo != nil && b.ElementTwo != nil && a.ElementTwo.Name == b.ElementTwo.Name))
}
