package models

import (
	"encoding/json"
	"os"
)

type Element struct {
	Name      string     `json:"name"`
	Recipes   [][]string `json:"recipes"`
	ImagePath string     `json:"image_path"`
}

func LoadElementsFromJSON(filePath string) ([]Element, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var elements []Element
	err = decoder.Decode(&elements)
	if err != nil {
		return nil, err
	}

	return elements, nil
}
