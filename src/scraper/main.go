package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	elements := scrapeElements(url)

	fmt.Printf("Total elemen ditemukan: %d\n", len(elements))
	elements = getMissingElementsIngredients(elements, url)
	fmt.Printf("Total elemen ditemukan: %d\n", len(elements))
	saveElementsToFile(elements, "../backend/data/elements.json")
}

func scrapeElements(url string) []Element {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	var wg sync.WaitGroup
	elementsChan := make(chan Element, 100) // Buffered channel to collect elements

	doc.Find("table.list-table.col-list.icon-hover").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				return // skip header
			}

			wg.Add(1)
			go func(row *goquery.Selection) {
				defer wg.Done()
				element := parseElement(row)
				if element != nil {
					elementsChan <- *element
				}
			}(row)
		})
	})

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(elementsChan)
	}()

	// Collect elements from the channel
	elements := []Element{}
	for element := range elementsChan {
		elements = append(elements, element)
	}

	return elements
}

func saveElementsToFile(elements []Element, filePath string) {
	log.Printf("Saving elements to file: %s", filePath)
	start := time.Now()

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Write JSON to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(elements); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("Elements saved to %s in %v", filePath, time.Since(start))
}

func getMissingElementsIngredients(elements []Element, url string) []Element {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch page for second pass: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status code error on second pass: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// 1. Collect all existing element names
	existing := make(map[string]bool)
	for _, e := range elements {
		existing[e.Name] = true
	}

	// 2. Collect all ingredient names from recipes
	used := make(map[string]bool)
	for _, e := range elements {
		for _, recipe := range e.Recipes {
			for _, ing := range recipe {
				used[ing] = true
			}
		}
	}

	// 3. Find missing ones
	var missing []string
	for ing := range used {
		if !existing[ing] {
			missing = append(missing, ing)
		}
	}
	log.Printf("Found %d missing ingredients. Attempting to scrape them...", len(missing))

	// 4. Try to scrape each missing element
	for _, name := range missing {
		row := findRowByElementName(doc, name)
		if row != nil {
			if el := parseElement(row); el != nil {
				elements = append(elements, *el)
			}
		} else {
			log.Printf("Could not find row for missing ingredient: %s", name)
			// Optional: Add placeholder
			imagePath := downloadImageFromIngredient(doc, name)
			elements = append(elements, Element{Name: name, Recipes: [][]string{}, ImagePath: imagePath})

		}
	}

	return elements
}

func findRowByElementName(doc *goquery.Document, name string) *goquery.Selection {
	var result *goquery.Selection
	doc.Find("table.list-table.col-list.icon-hover").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").EachWithBreak(func(i int, row *goquery.Selection) bool {
			text := row.Find("td").First().Text()
			if strings.EqualFold(strings.TrimSpace(text), name) {
				result = row
				return false
			}
			return true
		})
	})
	return result
}
