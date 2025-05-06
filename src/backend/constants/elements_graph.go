package constants

import "ccp/backend/models"

type Recipe struct {
	ElementOne models.Element `json:"element_one"`
	ElementTwo models.Element `json:"element_two"`
}

type ElementsGraphNode struct {
	Name            string               `json:"name"`
	ImagePath       string               `json:"image_path"`
	CanMakeElements []*ElementsGraphNode `json:"can_make_elements"`
	Recipes         []Recipe             `json:"recipes"`
}
