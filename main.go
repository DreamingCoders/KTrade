package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
)

// MarketplaceItem represents a single marketplace item
type MarketplaceItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Img   string `json:"img"`
	Link  string `json:"link"`
}

func fetchMarketplaceItems() ([]MarketplaceItem, error) {
	url := "https://kubeo.net/#/marketplace"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	items := []MarketplaceItem{}
	// Example selector matching your PHP logic; adjust if Kubeo changes layout
	doc.Find("div.grid a").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find("p.truncate").Text())
		if name == "" {
			name = "Unknown Item"
		}
		price := strings.TrimSpace(s.Find("p:contains('ph-currency-circle-dollar')").Text())
		if price == "" {
			price = "Free"
		}
		img, _ := s.Find("img").Attr("src")
		link, _ := s.Attr("href")
		if !strings.HasPrefix(link, "https://kubeo.net") {
			link = "https://kubeo.net" + link
		}
		items = append(items, MarketplaceItem{
			Name:  name,
			Price: price,
			Img:   img,
			Link:  link,
		})
	})
	return items, nil
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		items, err := fetchMarketplaceItems()
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Error fetching items: %v", err))
		}
		html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Kubeo Marketplace Mock API</title>
<style>
body { font-family: Arial, sans-serif; background: #121215; color: #eee; margin: 20px; }
.grid { display: flex; flex-wrap: wrap; gap: 15px; }
.card { background: #1b1b1f; padding: 10px; border-radius: 8px; width: 150px; text-align: center; transition: transform 0.2s; }
.card:hover { transform: scale(1.05); }
.card img { width: 100%; height: auto; border-radius: 4px; }
.card p { margin: 5px 0; }
.card a { color: #4fc3f7; text-decoration: none; }
</style>
</head>
<body>
<h1>Kubeo Marketplace (Mock API)</h1>
<p>This page demonstrates a mock "direct API" for fetching marketplace items.</p>
<div class="grid">`
		for _, item := range items {
			html += fmt.Sprintf(`<div class="card">
<a href="%s" target="_blank"><img src="%s" alt="%s"></a>
<p>%s</p>
<p><strong>%s</strong></p>
<a href="%s" target="_blank">View Item</a>
</div>`, item.Link, item.Img, item.Name, item.Name, item.Price, item.Link)
		}
		html += `</div></body></html>`
		return c.Type("html").SendString(html)
	})
	log.Println("Starting server on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
