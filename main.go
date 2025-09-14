package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
)

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

	// JSON API endpoint
	app.Get("/api/items", func(c *fiber.Ctx) error {
		items, err := fetchMarketplaceItems()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	})

	// SPA served directly
	app.Get("/", func(c *fiber.Ctx) error {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Kubeo Marketplace Mock SPA</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
<script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
</head>
<body class="bg-dark text-light">
<div id="app" class="container py-4">
  <h1 class="mb-4">Kubeo Marketplace</h1>
  <div class="mb-3">
    <button class="btn btn-primary me-2" @click="login">Login</button>
    <button class="btn btn-secondary" @click="register">Register</button>
  </div>
  <p v-if="user">Logged in as: {{ user }}</p>
  <div class="row row-cols-1 row-cols-md-4 g-4">
    <div class="col" v-for="item in items" :key="item.link">
      <div class="card h-100 bg-secondary text-light">
        <img :src="item.img" class="card-img-top" :alt="item.name">
        <div class="card-body">
          <h5 class="card-title">{{ item.name }}</h5>
          <p class="card-text"><strong>{{ item.price }}</strong></p>
          <a :href="item.link" target="_blank" class="btn btn-info btn-sm mb-2">View Item</a>
          <button class="btn btn-warning btn-sm" @click="sell(item)">Sell</button>
        </div>
      </div>
    </div>
  </div>
</div>
<script>
const { createApp } = Vue;
createApp({
  data() {
    return { items: [], user: null }
  },
  methods: {
    async loadItems() {
      try {
        const res = await fetch('/api/items');
        this.items = await res.json();
      } catch(e) { console.error(e) }
    },
    login() { this.user = prompt("Enter username:"); },
    register() { this.user = prompt("Choose username:"); },
    sell(item) {
      if(!this.user){ alert("Please login first."); return; }
      alert("Kubeo does not allow selling off-platform 😅");
      /* Uncomment below for mock sell code:
      console.log("Selling item:", item);
      alert("Mock sell: " + item.name + " for " + item.price);
      */
    }
  },
  mounted() { this.loadItems(); }
}).mount('#app');
</script>
</body>
</html>
		`
		return c.Type("html").SendString(html)
	})

	log.Println("Server running at http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
