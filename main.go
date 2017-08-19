package main

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func Crawl(url string, depth int, m *message) {
	defer func() { m.quit <- 0 }()

	// WebページからURLを取得
	urls, err := Fetch(url)

	// 結果送信
	m.res <- &respons{
		url: url,
		err: err,
	}

	if err == nil {
		for _, url := range urls {
			// 新しいリクエスト送信
			m.req <- &request{
				url:   url,
				depth: depth - 1,
			}
		}
	}
}

func Fetch(u string) (urls []string, err error) {
	baseUrl, err := url.Parse(u)
	if err != nil {
		return
	}

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	urls = make([]string, 0)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			reqUrl, err := baseUrl.Parse(href)
			if err == nil {
				urls = append(urls, reqUrl.String())
			}
		}
	})

	return
}

func main() {
	m := newMessage()
	go m.execute()
	m.req <- &request{
		url:   "https://www.google.co.jp/search?rlz=1C5CHFA_enJP693JP693&q=go+%E3%83%81%E3%83%A3%E3%83%8D%E3%83%AB&oq=go+%E3%83%81%E3%83%A3%E3%83%8D%E3%83%AB&gs_l=psy-ab.3...735206.735963.0.736286.9.6.0.0.0.0.221.221.2-1.1.0....0...1.1.64.psy-ab..9.0.0.donDWZ7eNLc",
		depth: 2,
	}
}
