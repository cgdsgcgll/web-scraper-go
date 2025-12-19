# Go Web Scraper

Bu proje, Go dili kullanılarak geliştirilmiş temel bir Web Scraper uygulamasıdır.  
Çalışma, Siber Tehdit İstihbaratı (CTI) kapsamında web sayfalarından veri toplama yeteneğini göstermek amacıyla hazırlanmıştır.

Projenin temel amacı; Go dili ile güvenilir HTTP iletişimi kurmak, web sayfalarının ham HTML içeriğini almak, ekran görüntüsü oluşturmak ve elde edilen verileri yerel dosyalara kaydetmektir.

---

## Özellikler

- Hedef URL’yi komut satırı argümanı olarak alır
- `net/http` kullanarak ham HTML içeriği çeker
- HTML çıktısını iki formatta kaydeder:
  - `.html` (web uyumlu)
  - `.txt` (Not Defteri ile kolay inceleme için)
- `chromedp` kullanarak **tam sayfa ekran görüntüsü** alır
- Sayfadaki tüm bağlantıları (`<a href="">`) çıkararak links.txt dosyasına kaydeder
- Aşağıdaki hataları yakalar ve kullanıcıya bildirir:
  - HTTP hataları (403, 404, 5xx)
  - DNS çözümleme hataları
  - Zaman aşımı (timeout) problemleri

---

## Kullanılan Teknolojiler

- **Go (Golang)**
- `net/http` – HTTP istekleri ve yanıt işleme
- `chromedp` – Headless Chrome ile ekran görüntüsü alma
- `goquery` – HTML parse etme ve link çıkarma

---

## Proje Yapısı

webscraper_go
├── main.go
├── go.mod
├── go.sum
├── README.md
└── results/
├── www.nasa.gov*/
│ ├── site*data.html
│ ├── site_data.txt
│ ├── screenshot.png
│ └── links.txt
├── www.wikipedia.org*/
├── github.com*trending/
├── arxiv.org_list_cs_recent/
├── home.cern*/
├── pastebin.com*archive/
├── support.torproject.org*/
├── www.ietf.org*/
├── www.kutahya.bel.tr*/
├── www.saglik.gov.tr*/
├── www.tccb.gov.tr*/
├── www.tcmb.gov.tr*wps_wcm_connect_TR_TCMB+TR_Main+Menu_Yayinlar/
├── www.tubitak.gov.tr*/
├── www.turkiye.gov.tr*/
└── www.yok.gov.tr*/

`results/` klasörü altında, test edilen her web sitesi için ayrı bir klasör oluşturulur.  
Bu klasörlerin içinde:

- `site_data.html` → Web sayfasının ham HTML içeriği
- `site_data.txt` → Not Defteri ile okunabilir HTML içeriği
- `screenshot.png` → Sayfanın tam ekran görüntüsü
- `links.txt` → Sayfa içerisindeki bağlantılar

bulunmaktadır.

---

## Kullanım

Programı çalıştırmak için proje dizininde aşağıdaki komut kullanılır:

```bash
go run main.go https://www.nasa.gov/

```
