package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Name      string     `json:"name"`
	Recipes   [][]string `json:"recipes"`
	ImagePath string     `json:"image_path"`
}

func parseElement(row *goquery.Selection) *Element {
	cells := row.Find("td")
	if cells.Length() != 2 {
		return nil
	}

	elementName := strings.TrimSpace(cells.Eq(0).Text())
	if elementName == "" || strings.ToLower(elementName) == "element" {
		return nil
	}

	recipes := parseRecipes(cells.Eq(1))
	imagePath := downloadImage(cells.Eq(0), elementName)

	return &Element{
		Name:      elementName,
		Recipes:   recipes,
		ImagePath: imagePath,
	}
}

func parseRecipes(cell *goquery.Selection) [][]string {
	recipes := [][]string{}
	cell.Find("li").Each(func(_ int, li *goquery.Selection) {
		ingredients := []string{}
		li.Find("a").Each(func(_ int, a *goquery.Selection) {
			text := strings.TrimSpace(a.Text())
			if text != "" && strings.ToLower(text) != "file" {
				ingredients = append(ingredients, text)
			}
		})
		if len(ingredients) > 0 {
			recipes = append(recipes, ingredients)
		}
	})

	return recipes
}
