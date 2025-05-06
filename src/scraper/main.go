package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	elements := scrapeElements(url)

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

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(elements); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	log.Printf("Elements saved to %s in %v", filePath, time.Since(start))
}
