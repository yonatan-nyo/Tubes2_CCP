package models

import (
	"fmt"
	"strings"
)

// Updated Debug function with maxDepth parameter
func Debug(root *ElementsGraphNode, maxDepth int) {
	DebugBasicElementsFromRoot()
	fmt.Println("\n=== Elements Graph Debug Output ===")
	visited := make(map[string]bool)
	for _, recipe := range ElementsGraph.RecipesToMakeOtherElement {
		printNodeWithMaxDepth(recipe.ElementOne, visited, 0, maxDepth)
	}
}

// Original Debug function for backward compatibility
func DebugDefault(root *ElementsGraphNode) {
	Debug(root, 1) // Default to a depth of 1 to prevent too much output
}

func DebugBasicElementsFromRoot() {
	fmt.Println("=== Basic Elements from Root Node ===")
	for _, recipe := range ElementsGraph.RecipesToMakeOtherElement {
		if recipe.ElementTwo == nil {
			fmt.Printf("- %s (%s)\n", recipe.ElementOne.Name, recipe.ElementOne.ImagePath)
		}
	}
}

func printNodeWithMaxDepth(node *ElementsGraphNode, visited map[string]bool, depth int, maxDepth int) {
	if node == nil || visited[node.Name] || (maxDepth >= 0 && depth > maxDepth) {
		return
	}
	visited[node.Name] = true

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s- (%d) Element: %s (%s)\n", indent, node.Tier, node.Name, node.ImagePath)

	// Print recipes to make this element
	if len(node.RecipesToMakeThisElement) > 0 {
		fmt.Printf("%s  Recipes to make this: (%d)\n", indent, len(node.RecipesToMakeThisElement))
		for _, r := range node.RecipesToMakeThisElement {
			if r.ElementTwo != nil {
				fmt.Printf("%s    %s(%d) + %s(%d) => %s\n", indent, r.ElementOne.Name, r.ElementOne.Tier, r.ElementTwo.Name, r.ElementTwo.Tier, node.Name)
			} else {
				fmt.Printf("%s    %s => %s\n", indent, r.ElementOne.Name, node.Name)
			}
		}
	}

	// Print recipes where this element is used to make others
	// if len(node.RecipesToMakeOtherElement) > 0 {
	// 	fmt.Printf("%s  Recipes using this element to make others: (%d)\n", indent, len(node.RecipesToMakeOtherElement))
	// 	for _, r := range node.RecipesToMakeOtherElement {
	// 		// We need to determine what's the other element in the recipe and what's the target
	// 		if node == r.ElementOne && r.ElementTwo != nil {
	// 			fmt.Printf("%s    %s + %s => %s\n", indent, node.Name, r.ElementTwo.Name, r.TargetElementName)
	// 		} else if node == r.ElementTwo && r.ElementOne != nil {
	// 			fmt.Printf("%s    %s + %s => %s\n", indent, node.Name, r.ElementOne.Name, r.TargetElementName)
	// 		} else if r.ElementTwo == nil {
	// 			// This is a basic element case
	// 			fmt.Printf("%s    %s => %s\n", indent, node.Name, r.TargetElementName)
	// 		}
	// 	}
	// }

	// Recursively print connected nodes
	for _, r := range node.RecipesToMakeOtherElement {
		var target *ElementsGraphNode

		// Find the target element (not the current node)
		if r.ElementOne != node && r.ElementOne != nil {
			target = r.ElementOne
		} else if r.ElementTwo != nil {
			target = r.ElementTwo
		}

		// If we found a valid target and it's not the current node, recursively print it
		if target != nil && target != node {
			printNodeWithMaxDepth(target, visited, depth+1, maxDepth)
		}
	}
}

// Add a convenient debug function that lets you debug a specific element
func DebugElement(elementName string, maxDepth int) {
	fmt.Printf("\n=== Debug for Element: %s ===\n", elementName)
	node, exists := nameToNode[elementName]
	if !exists {
		fmt.Printf("Element '%s' not found in the graph.\n", elementName)
		return
	}

	visited := make(map[string]bool)
	printNodeWithMaxDepth(node, visited, 0, maxDepth)
}
