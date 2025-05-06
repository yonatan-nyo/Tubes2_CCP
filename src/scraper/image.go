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
