package models

import "fmt"

type SafeRecipe struct {
	ElementOneName    string `json:"element_one"`
	ElementTwoName    string `json:"element_two,omitempty"` // omit if nil
	TargetElementName string `json:"target_element_name"`
}

type SafeElementNode struct {
	Name                      string       `json:"name"`
	ImagePath                 string       `json:"image_path"`
	RecipesToMakeThisElement  []SafeRecipe `json:"recipes_to_make_this_element"`
	RecipesToMakeOtherElement []SafeRecipe `json:"recipes_to_make_other_element"`
}

func ToSafeGraph(root *ElementsGraphNode) []SafeElementNode {
	visited := make(map[string]bool)
	safeNodes := make([]SafeElementNode, 0)

	var dfs func(node *ElementsGraphNode)
	dfs = func(node *ElementsGraphNode) {
		if node == nil || visited[node.Name] {
			return
		}
		visited[node.Name] = true

		safeNode := SafeElementNode{
			Name:      node.Name,
			ImagePath: node.ImagePath,
		}

		for _, r := range node.RecipesToMakeThisElement {
			safeNode.RecipesToMakeThisElement = append(safeNode.RecipesToMakeThisElement, SafeRecipe{
				ElementOneName:    r.ElementOne.Name,
				ElementTwoName:    safeName(r.ElementTwo),
				TargetElementName: r.TargetElementName,
			})
			// Traverse ingredients
			dfs(r.ElementOne)
			if r.ElementTwo != nil {
				dfs(r.ElementTwo)
			}
		}

		for _, r := range node.RecipesToMakeOtherElement {
			safeNode.RecipesToMakeOtherElement = append(safeNode.RecipesToMakeOtherElement, SafeRecipe{
				ElementOneName:    r.ElementOne.Name,
				ElementTwoName:    safeName(r.ElementTwo),
				TargetElementName: r.TargetElementName,
			})
			fmt.Println("Adding recipe to make other element:", r.ElementOne.Name, r.ElementTwo, r.TargetElementName)
			// Traverse result element
			if target, ok := nameToNode[r.TargetElementName]; ok {
				dfs(target)
			}
		}

		safeNodes = append(safeNodes, safeNode)
	}

	dfs(root)
	return safeNodes
}

func safeName(node *ElementsGraphNode) string {
	if node == nil {
		return ""
	}
	return node.Name
}
