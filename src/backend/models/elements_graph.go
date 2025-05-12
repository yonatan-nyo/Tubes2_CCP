package models

import (
	"os"
	"strings"
)

type ElementsGraphNode struct {
	Name                      string    	  `json:"name"`
	ImagePath                 string    	  `json:"image_path"`
	RecipesToMakeThisElement  []*Recipe 	  `json:"recipes_to_make_this_element"`
	RecipesToMakeOtherElement []*Recipe 	  `json:"recipes_to_make_other_element"`
	Tier                      int       	  `json:"tier"`
	IsVisited                 bool      	  `json:"is_visited"`
	MadeFrom				  map[string]bool `json:"made_from"`
}

type Recipe struct {
	ElementOne        *ElementsGraphNode `json:"element_one"`
	ElementTwo        *ElementsGraphNode `json:"element_two"`
	TargetElementName string             `json:"target_element_name"`
}

var ElementsGraph = &ElementsGraphNode{
	Name:                      "Root",
	RecipesToMakeThisElement:  []*Recipe{},
	RecipesToMakeOtherElement: []*Recipe{},
	IsVisited:                 false,
}

var nameToNode = make(map[string]*ElementsGraphNode)

func GetElementsGraphNodeByName(name string) (*ElementsGraphNode, bool) {
	node, exists := nameToNode[name]
	return node, exists
}

func (node *ElementsGraphNode) IsThisMadeFrom(element string) bool {
	if node.MadeFrom == nil {
		return false
	}
	return node.MadeFrom[element]
}

type ElementsGraphNodeDTO struct {
	Name                      string      `json:"name"`
	ImagePath                 string      `json:"image_path"`
	RecipesToMakeThisElement  []RecipeDTO `json:"recipes_to_make_this_element"`
	RecipesToMakeOtherElement []RecipeDTO `json:"recipes_to_make_other_element"`
	IsVisited                 bool        `json:"is_visited"`
	Tier                      int         `json:"tier"`
}

type RecipeDTO struct {
	ElementOneName    string `json:"element_one"`
	ElementTwoName    string `json:"element_two"`
	TargetElementName string `json:"target_element_name"`
}

func GetJSONDTONodes() []ElementsGraphNodeDTO {
	nameToNodeList := make([]ElementsGraphNodeDTO, 0)
	for _, node := range nameToNode {
		dto := ElementsGraphNodeDTO{
			Name:                      node.Name,
			ImagePath:                 GetImagePath(node.ImagePath),
			RecipesToMakeThisElement:  make([]RecipeDTO, len(node.RecipesToMakeThisElement)),
			RecipesToMakeOtherElement: make([]RecipeDTO, len(node.RecipesToMakeOtherElement)),
			IsVisited:                 node.IsVisited,
			Tier:                      node.Tier,
		}

		for i, recipe := range node.RecipesToMakeThisElement {
			dto.RecipesToMakeThisElement[i] = RecipeDTO{
				ElementOneName:    recipe.ElementOne.Name,
				ElementTwoName:    recipe.ElementTwo.Name,
				TargetElementName: recipe.TargetElementName,
			}
		}

		for i, recipe := range node.RecipesToMakeOtherElement {
			dto.RecipesToMakeOtherElement[i] = RecipeDTO{
				ElementOneName:    recipe.ElementOne.Name,
				ElementTwoName:    recipe.ElementTwo.Name,
				TargetElementName: recipe.TargetElementName,
			}
		}

		nameToNodeList = append(nameToNodeList, dto)
	}
	return nameToNodeList
}

func GetImagePath(path string) string {
	if path == "" {
		return ""
	}
	path = strings.TrimPrefix(path, "../backend/")
	path = strings.TrimPrefix(path, "/")
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://nyo.kirisame.jp.net/"
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	// Concatenate the base URL with the normalized path
	return baseURL + path
}

func GetElementsFromNameToNodeDTO() []*ElementsGraphNodeDTO {
	// change recipe to string from nameToNode
	nameToNodeList := make([]*ElementsGraphNodeDTO, 0)
	for _, node := range nameToNode {
		dto := &ElementsGraphNodeDTO{
			Name:                      node.Name,
			ImagePath:                 GetImagePath(node.ImagePath),
			RecipesToMakeThisElement:  make([]RecipeDTO, len(node.RecipesToMakeThisElement)),
			RecipesToMakeOtherElement: make([]RecipeDTO, len(node.RecipesToMakeOtherElement)),
			IsVisited:                 node.IsVisited,
			Tier:                      node.Tier,
		}

		for i, recipe := range node.RecipesToMakeThisElement {
			dto.RecipesToMakeThisElement[i] = RecipeDTO{
				ElementOneName:    recipe.ElementOne.Name,
				ElementTwoName:    recipe.ElementTwo.Name,
				TargetElementName: recipe.TargetElementName,
			}
		}

		for i, recipe := range node.RecipesToMakeOtherElement {
			dto.RecipesToMakeOtherElement[i] = RecipeDTO{
				ElementOneName:    recipe.ElementOne.Name,
				ElementTwoName:    recipe.ElementTwo.Name,
				TargetElementName: recipe.TargetElementName,
			}
		}

		nameToNodeList = append(nameToNodeList, dto)
	}
	return nameToNodeList
}
