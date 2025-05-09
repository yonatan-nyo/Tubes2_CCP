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
	cell.Find("img").Each(func(_ int, img *goquery.Selection) {
		src, exists := img.Attr("data-src")
		if !exists {
			src, exists = img.Attr("src")
		}
		if exists {
			imagePath = src

			encodedName := url.QueryEscape(strings.ReplaceAll(elementName, " ", "_"))
			imageFilePath := fmt.Sprintf("../backend/public/%s.png", encodedName)

			// Check if the file already exists
			if _, err := os.Stat(imageFilePath); err == nil {
				log.Printf("Image for %s already exists at %s, skipping download", elementName, imageFilePath)
				imagePath = imageFilePath
				return
			}

			log.Printf("\nDownloading image for %s", elementName)
			start := time.Now()

			resp, err := http.Get(imagePath)
			if err != nil {
				log.Printf("Failed to download image for %s: %v", elementName, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Printf("Failed to download image for %s: status code %d", elementName, resp.StatusCode)
				return
			}

			file, err := os.Create(imageFilePath)
			if err != nil {
				log.Printf("Failed to create image file for %s: %v", elementName, err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Printf("Failed to save image for %s: %v", elementName, err)
				return
			}

			log.Printf("Image for %s downloaded and saved to %s in %v", elementName, imageFilePath, time.Since(start))
			imagePath = imageFilePath
		}
	})

	return imagePath
}

func downloadImageFromIngredient(doc *goquery.Document, ingredientName string) string {
	var imagePath string

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
							encoded := url.QueryEscape(strings.ReplaceAll(ingredientName, " ", "_"))
							imageFilePath := fmt.Sprintf("../backend/public/%s.png", encoded)
							if _, err := os.Stat(imageFilePath); err == nil {
								return
							}

							resp, err := http.Get(src)
							if err != nil {
								log.Printf("Failed to download ingredient image for %s: %v", ingredientName, err)
								return
							}
							defer resp.Body.Close()

							file, err := os.Create(imageFilePath)
							if err != nil {
								log.Printf("Failed to create image file for %s: %v", ingredientName, err)
								return
							}
							defer file.Close()

							_, err = io.Copy(file, resp.Body)
							if err != nil {
								log.Printf("Failed to save ingredient image for %s: %v", ingredientName, err)
								return
							}
							log.Printf("Downloaded image for ingredient %s", ingredientName)
							imagePath = imageFilePath
						}
					}
				}
			})
		})
	})

	return imagePath
}
