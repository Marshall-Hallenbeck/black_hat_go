package main

import (
	"archive/zip"
	"black_hat_go/3_http_and_remote_interaction/bing_scraping/metadata"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// accepts a goquery.Selection instance, which will be populated with an anchor HTML element
func handler(i int, s *goquery.Selection) {
	// find and extract the href attribute
	url, ok := s.Find("a").Attr("href")
	if !ok {
		return
	}
	// get the URL extracted from the href
	fmt.Printf("%d: %s\n", i, url)
	res, err := http.Get(url)
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return
	}

	cp, ap, err := metadata.NewProperties(r)
	if err != nil {
		return
	}

	log.Printf(
		"%21s %s - %s %s\n",
		cp.Creator,
		cp.LastModifiedBy,
		ap.Application,
		ap.GetMajorVersion(),
	)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Missing required argument. Usage: main.go domain ext")
	}
	domain := os.Args[1]
	// good filetypes: docx, xlsx, pptx
	// some macro alternatives: docm, xlsm, pptm
	filetype := os.Args[2]

	client := &http.Client{}
	q := fmt.Sprintf(
		"site:%s && filetype:%s && instreamset:(url title):%s",
		domain,
		filetype,
		filetype,
	)
	search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))
	req, err := http.NewRequest("GET", search, nil)
	if err != nil {
		return
	}

	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error!")
		log.Panicln(err)
	}
	defer res.Body.Close()
	s := "html body div#b_content ol#b_results li.b_algo h2"
	doc.Find(s).Each(handler)
}
