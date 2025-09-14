<?php
/**
 * Kubeo Marketplace Mock API
 * --------------------------
 * This script demonstrates fetching items from Kubeo Marketplace.
 * It does not rely on a real API, but instead scrapes the HTML page.
 * It also formats items in a way similar to trading sites like rbx.trade.
 *
 * Features:
 *  - Scrape item name, price, and image
 *  - Generate direct links to each item page (e.g., https://kubeo.net/#/item/136)
 *  - Mock “direct API” endpoint style for demonstration
 *
 * NOTE: This is a mock implementation and intended for showcasing purposes only.
 */
$marketplaceURL = "https://kubeo.net/#/marketplace"; // URL of the marketplace page (SPA hash route; may require headless JS in production)
// cURL request to fetch page HTML. This version basically grabs it via HTML there is an API version as well depending on what you want to really do honestly.
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $marketplaceURL);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
curl_setopt($ch, CURLOPT_USERAGENT, "Mozilla/5.0"); // mimic browser
$html = curl_exec($ch);
curl_close($ch);
// Use DOMDocument to parse HTML
libxml_use_internal_errors(true);
$dom = new DOMDocument();
$dom->loadHTML($html);
libxml_clear_errors();
$xpath = new DOMXPath($dom);
// Extract items from the page
$items = $xpath->query("//div[contains(@class,'grid')]//a");
$mockAPI = [];
foreach ($items as $item) {
    $name = $xpath->query(".//p[contains(@class,'truncate')]", $item)->item(0)?->nodeValue ?? "Unknown Item";
    $priceNode = $xpath->query(".//p[contains(., 'ph-currency-circle-dollar')]", $item)->item(0);
    $price = $priceNode ? trim(str_replace(["\n", "\t"], "", $priceNode->nodeValue)) : "Free";
    $img = $xpath->query(".//img[contains(@src,'items')]", $item)->item(0)?->getAttribute('src') ?? "";
    $link = $item->getAttribute('href') ?: "#";
    // Push to mock API array
    $mockAPI[] = [
        "name" => $name,
        "price" => $price,
        "img" => $img,
        "link" => "https://kubeo.net" . $link // direct item page
    ];
}
// Render items in a “trade-style” layout
?>
<!DOCTYPE html>
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
<div class="grid">
<?php foreach ($mockAPI as $item): ?>
    <div class="card">
        <a href="<?= $item['link'] ?>" target="_blank">
            <img src="<?= $item['img'] ?>" alt="<?= htmlspecialchars($item['name']) ?>">
        </a>
        <p><?= htmlspecialchars($item['name']) ?></p>
        <p><strong><?= $item['price'] ?></strong></p>
        <a href="<?= $item['link'] ?>" target="_blank">View Item</a>
    </div>
<?php endforeach; ?>
</div>
<script>
// Mock JavaScript logic for tabbed interface (like a trade site dashboard)
document.addEventListener('DOMContentLoaded', () => {
    const tabs = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');
    tabs.forEach(tab => {
        tab.addEventListener('click', (e) => {
            e.preventDefault();
            tabs.forEach(t => t.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));
            const tabId = tab.dataset.tab;
            const content = document.getElementById(tabId);
            tab.classList.add('active');
            if (content) content.classList.add('active');
        });
    });
});
</script>
</body>
</html>
