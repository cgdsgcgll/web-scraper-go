package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go <URL>")
		fmt.Println("Örnek  : go run main.go https://example.com")
		os.Exit(1)
	}
	url := os.Args[1]

	// Klasör adı ürettiğim kısım
	safe := strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
	safe = strings.ReplaceAll(safe, "/", "_")
	outDir := "results\\" + safe

	_ = os.MkdirAll(outDir, 0755)

	// Çıktı dosya yolları
	// Fazladan site_data olmasının sebebi birinin html diğerinin txt uzantılı olması
	htmlPath := outDir + "\\site_data.html"
	txtPath := outDir + "\\site_data.txt"
	ssPath := outDir + "\\screenshot.png"
	linksPath := outDir + "\\links.txt"

	// HTML çekilen kısım
	html, code, err := fetchHTML(url)
	fmt.Printf("HTTP Durum Kodu: %d\n", code)

	if err != nil {
		fmt.Printf("Hata: %v\n", err)

		// Hata olsa bile dönen body varsa kaydedilen kısım
		if html != "" {
			_ = os.WriteFile(htmlPath, []byte(html), 0644)
			_ = os.WriteFile(txtPath, []byte(html), 0644)
			fmt.Println("Yine de dönen içerik kaydedildi:", htmlPath, "ve", txtPath)
		}
		os.Exit(1)
	}

	// HTML'i hem .html hem .txt olarak kaydettiğim kısm
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		fmt.Printf("HTML kaydetme hatası: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(txtPath, []byte(html), 0644); err != nil {
		fmt.Printf("TXT kaydetme hatası: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("HTML kaydedildi:", htmlPath)
	fmt.Println("TXT kaydedildi :", txtPath)

	// Screenshot alınan kısım
	if err := takeScreenshot(url, ssPath); err != nil {
		fmt.Printf("Screenshot alma hatası: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Screenshot kaydedildi:", ssPath)

	// linkleri çıkardığım kısım
	if err := extractLinks(html, linksPath); err != nil {
		fmt.Printf("Link çıkarma hatası (bonus): %v\n", err)
	} else {
		fmt.Println("Linkler kaydedildi:", linksPath)
	}

	fmt.Println("✅ Tamamlandı.")
}

func fetchHTML(url string) (string, int, error) {
	client := &http.Client{Timeout: 25 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("istek oluşturulamadı: %w", err)
	}

	// Bazı siteler bot engeli uyguladığından dolayı header'ları daha gerçekçi yaptım.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://www.google.com/")

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("bağlantı hatası: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("HTML okunamadı: %w", err)
	}

	// 2xx değilse hata olarak döndür 
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return string(body), resp.StatusCode, fmt.Errorf("HTTP hata kodu: %d (%s)", resp.StatusCode, resp.Status)
	}

	return string(body), resp.StatusCode, nil
}

func takeScreenshot(url, outFile string) error {
	// Daha stabil başlatmak için allocator ayarları
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Ağır siteler için timeout artırdığım kısım
	ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return fmt.Errorf("chromedp run hatası: %w", err)
	}

	return os.WriteFile(outFile, buf, 0644)
}

func extractLinks(html, outFile string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return fmt.Errorf("parse hatası: %w", err)
	}

	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("links dosyası oluşturulamadı: %w", err)
	}
	defer f.Close()

	seen := map[string]bool{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" && !seen[href] {
			seen[href] = true
			fmt.Fprintln(f, href)
		}
	})
	return nil
}
