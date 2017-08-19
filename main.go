package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"

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

	// 取得したhtmlを文字列で確認したい時はこれ
	//body, err := ioutil.ReadAll(resp.Body)
	//buf := bytes.NewBuffer(body)
	//html := buf.String()
	//fmt.Println(html)
	// ---------------

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	urls = make([]string, 0)
	doc.Find(".r").Each(func(_ int, srg *goquery.Selection) {
		srg.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				reqUrl, err := baseUrl.Parse(href)
				if err == nil {
					urls = append(urls, reqUrl.String())
				}
			}
		})
	})

	return
}

func main() {
	var word = flag.String("w", " ", "検索ワードを入力して下さい")
	flag.Parse()
	log.Println("検索ワード：", *word)
	*word = strings.Replace(*word, " ", "+", -1)
	firstURL := "https://www.google.co.jp/search?rlz=1C5CHFA_enJP693JP693&q=" + string(*word)
	log.Println("検索URL：", firstURL)
	m := newMessage()
	go m.execute()
	m.req <- &request{
		url:   firstURL,
		depth: 2,
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndSearver:", err)
	}
}
