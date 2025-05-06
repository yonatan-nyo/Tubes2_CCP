package models

type ElementsGraphNode struct {
	Name                      string    `json:"name"`
	ImagePath                 string    `json:"image_path"`
	RecipesToMakeThisElement  []*Recipe `json:"recipes_to_make_this_element"`
	RecipesToMakeOtherElement []*Recipe `json:"recipes_to_make_other_element"`
	IsVisited                 bool      `json:"is_visited"`
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

type RecipeTreeNode struct {
	Name      string          `json:"name"`
	ImagePath string          `json:"image_path"`
	Recipe    *Recipe         `json:"recipe"`
	Child     *RecipeTreeNode `json:"child"`
}

// RecipeTreeNode can point to the same child indicating that they are the parent(recipe) of the same child
