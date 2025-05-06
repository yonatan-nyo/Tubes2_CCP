package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func main() {
    url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

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

    elementRecipes := make(map[string][][]string)

    doc.Find("table.list-table.col-list.icon-hover").Each(func(_ int, table *goquery.Selection) {
        table.Find("tr").Each(func(i int, row *goquery.Selection) {
            if i == 0 {
                return // skip header
            }

            cells := row.Find("td")
            if cells.Length() != 2 {
                return
            }

            elementName := strings.TrimSpace(cells.Eq(0).Text())
            if elementName == "" || strings.ToLower(elementName) == "element" {
                return
            }

            recipes := [][]string{}
            cells.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
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

            if len(recipes) == 0 {
                recipes = append(recipes, []string{}) // element dasar
            }

            elementRecipes[elementName] = recipes
        })
    })

    fmt.Printf("Total elemen ditemukan: %d\n", len(elementRecipes))
    count := 0
    for name, recs := range elementRecipes {
        fmt.Printf("%s:\n", name)
        for _, r := range recs {
            fmt.Printf("  - %v\n", r)
        }
        count++
        if count >= 5 {
            break
        }
    }

    // Simpan ke file JSON
    file, err := os.Create("little_alchemy2_recipes.json")
    if err != nil {
        log.Fatalf("Failed to create file: %v", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(elementRecipes); err != nil {
        log.Fatalf("Failed to write JSON: %v", err)
    }
}
