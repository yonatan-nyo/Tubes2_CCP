package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func downloadImage(cell *goquery.Selection, elementName string) string {
	imagePath := ""

	// Check for image source in the `img` tag
	cell.Find("img").Each(func(_ int, img *goquery.Selection) {
		src, exists := img.Attr("data-src")
		if !exists {
			src, exists = img.Attr("src")
		}
		if exists {
			// Encode the element name for the image file
			encodedName := strings.ReplaceAll(elementName, " ", "_")
			imageFilePath := fmt.Sprintf("../backend/public/%s.png", encodedName)

			// Ensure the directory exists
			err := os.MkdirAll("../backend/public", os.ModePerm)
			if err != nil {
				log.Printf("Failed to create directories for %s: %v", elementName, err)
				return
			}

			// Check if the file already exists
			if _, err := os.Stat(imageFilePath); err == nil {
				log.Printf("Image for %s already exists at %s, skipping download", elementName, imageFilePath)
				imagePath = imageFilePath
				return
			}

			// Download the image
			log.Printf("\nDownloading image for %s", elementName)
			start := time.Now()

			resp, err := http.Get(src)
			if err != nil {
				log.Printf("Failed to download image for %s: %v", elementName, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Printf("Failed to download image for %s: status code %d", elementName, resp.StatusCode)
				return
			}

			// Create the image file
			file, err := os.Create(imageFilePath)
			if err != nil {
				log.Printf("Failed to create image file for %s: %v", elementName, err)
				return
			}
			defer file.Close()

			// Save the image content to the file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Printf("Failed to save image for %s: %v", elementName, err)
				return
			}

			log.Printf("Image for %s downloaded and saved to %s in %v", elementName, imageFilePath, time.Since(start))

			// Return the relative URL path of the image
			imagePath = "/public/" + url.PathEscape(encodedName+".png")
		}
	})

	return imagePath
}

func downloadImageFromIngredient(doc *goquery.Document, ingredientName string) string {
	var imagePath string

	// Find the table containing ingredient rows
	doc.Find("table.list-table.col-list.icon-hover").Each(func(_ int, table *goquery.Selection) {
		table.Find("td").Each(func(_ int, cell *goquery.Selection) {
			cell.Find("a").Each(func(_ int, a *goquery.Selection) {
				if strings.EqualFold(strings.TrimSpace(a.Text()), ingredientName) {
					img := a.Find("img")
					if img.Length() == 0 {
						img = a.Parent().Find("img")
					}
					if img.Length() > 0 {
						src, exists := img.Attr("data-src")
						if !exists {
							src, exists = img.Attr("src")
						}
						if exists {
							// Use raw (not escaped) name to save
							filename := strings.ReplaceAll(ingredientName, " ", "_") + ".png"
							rawFilePath := fmt.Sprintf("../backend/public/%s", filename)

							// Ensure the directory exists
							err := os.MkdirAll("../backend/public", os.ModePerm)
							if err != nil {
								log.Printf("Failed to create directories for %s: %v", ingredientName, err)
								return
							}

							// Check if the file already exists
							if _, err := os.Stat(rawFilePath); err == nil {
								// If file exists, return the relative path
								imagePath = "/public/" + url.PathEscape(filename)
								return
							}

							// Download and save the image if it doesn't exist
							resp, err := http.Get(src)
							if err != nil {
								log.Printf("Failed to download image for %s: %v", ingredientName, err)
								return
							}
							defer resp.Body.Close()

							// Create the image file
							file, err := os.Create(rawFilePath)
							if err != nil {
								log.Printf("Failed to create file for %s: %v", ingredientName, err)
								return
							}
							defer file.Close()

							// Save the image content to the file
							_, err = io.Copy(file, resp.Body)
							if err != nil {
								log.Printf("Failed to save image for %s: %v", ingredientName, err)
								return
							}

							log.Printf("Downloaded image for ingredient: %s", ingredientName)

							// Return the relative URL path of the image
							imagePath = "/public/" + url.PathEscape(filename)
						}
					}
				}
			})
		})
	})

	return imagePath
}
