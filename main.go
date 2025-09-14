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

	app.Get("/api/items", func(c *fiber.Ctx) error {
		items, err := fetchMarketplaceItems()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(items)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Kubeo SPA</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
<script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
</head>
<body class="bg-dark text-light">
<div id="app" class="container py-4">
  <h1 class="mb-4">Kubeo SPA</h1>
  <div class="mb-3">
    <button class="btn btn-primary me-2" @click="login">Login</button>
    <button class="btn btn-secondary me-2" @click="register">Register</button>
    <button class="btn btn-info me-2" @click="navigate('marketplace')">Marketplace</button>
    <button class="btn btn-success me-2" @click="navigate('forums')">Forums</button>
    <button class="btn btn-warning me-2" @click="navigate('leaderboard')">Leaderboard</button>
    <button class="btn btn-light text-dark" @click="navigate('profile')">Profile</button>
  </div>
  <p v-if="user">Logged in as: {{ user }}</p>

  <!-- Marketplace Page -->
  <div v-if="currentPage==='marketplace'" class="row row-cols-1 row-cols-md-4 g-4">
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

  <!-- Forums Page -->
  <div v-if="currentPage==='forums'">
    <h2>Forums</h2>
    <p>Coming soon: Discuss Kubeo items and games here!</p>
    <div v-for="n in 5" class="card bg-secondary text-light mb-2 p-2">Forum post #{{ n }}</div>
  </div>

  <!-- Leaderboard Page -->
  <div v-if="currentPage==='leaderboard'">
    <h2>Leaderboard</h2>
    <ol class="list-group list-group-numbered">
      <li class="list-group-item bg-secondary text-light">Player1 - 1500 pts</li>
      <li class="list-group-item bg-secondary text-light">Player2 - 1200 pts</li>
      <li class="list-group-item bg-secondary text-light">Player3 - 1100 pts</li>
    </ol>
  </div>

  <!-- Profile Page -->
  <div v-if="currentPage==='profile'">
    <h2>Profile</h2>
    <p v-if="!user">Please login to see profile info.</p>
    <div v-else>
      <p>Username: {{ user }}</p>
      <p>Joined: 2025-01-01</p>
      <p>Items owned: 5</p>
    </div>
  </div>
</div>

<script>
const { createApp } = Vue;
createApp({
  data() { return { items: [], user: null, currentPage: 'marketplace' } },
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
    },
    navigate(page) { this.currentPage = page; }
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
