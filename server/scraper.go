package main 

import (
	"log"
	"net/http"

	"strings"
	"regexp"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/DavidBelicza/TextRank/v2"
)

type Page struct {
	Url 	string
	Summary []string
	Links 	[]string
}

func getLinks(doc *goquery.Document, scraped *sync.Map, base_url string) []string {
	var res []string
	traversed := make(map[string]bool)
	doc.Find(".mw-parser-output a").Each(func(i int, s *goquery.Selection) {
		link_raw, exists := s.Attr("href")
		if exists {
			link, err := url.Parse(link_raw)	
			if err != nil {
				return
			}

			if link.String() == "" || 
				strings.Contains(link.String(), ":") ||
				strings.Contains(link.String(), "Index") ||
				strings.Contains(link.String(), "List") ||
				link.RawQuery != "" ||
				link.Fragment != "" ||
				link.Host 	  != "" {
				return
			}
			
			_, scrp := scraped.Load(base_url + link.String())
			if !scrp && 
				!traversed[link.String()] && 
				len(res) < 5 {

				res = append(res, base_url + link.String())
				traversed[link.String()] = true
			}
		}
	})
	return res
}

var rule = textrank.NewDefaultRule()
var lang = textrank.NewDefaultLanguage()
var algorithm = textrank.NewDefaultAlgorithm()
// matches all '[\d]' or any newlines ('\n')
var reg = regexp.MustCompile(`(\[[0-9]+\])|(\n)`)

func getSummary(doc *goquery.Document) []string {
	var text string
	doc.Find(".mw-parser-output p").Each(func(i int, s *goquery.Selection) {
		text += s.Text();
	})

	tr := textrank.NewTextRank() 
	tr.Populate(text, lang, rule)
	sentences := textrank.FindSentencesByRelationWeight(tr, 5)

	var res []string
	for i, s := range sentences {
		if i >= 5 { break }	
		res = append(res, reg.ReplaceAllString(s.Value, ""))
	}

	return res
}

func scrape(page string, depth int, scraped *sync.Map, wg *sync.WaitGroup, result *[]Page) {
	if depth == 0 { return }

	_, contains := scraped.Load(page)
	if contains {
		return 
	}
	scraped.Store(page, true)

	base_url, err := url.Parse(page)
	if err != nil {
		return 
	}
	base_url.RawQuery = ""
	base_url.Fragment = ""

	// check if we have already traversed this page
		
	// get the page html
	res, err := http.Get(page)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
		return 
	}

	// load the html query document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var curr Page
	curr.Url = base_url.String()
	base_url.Path = ""
	curr.Links = getLinks(doc, scraped, base_url.String())
	curr.Summary = getSummary(doc)

	if len(curr.Links) > 5 { curr.Links = curr.Links[:5] }
	*result = append(*result, curr)

	for _, l := range curr.Links {
		wg.Add(1)
		go func(l string) {
			scrape(l, depth-1, scraped, wg, result)
			wg.Done()
		}(l)
	}

}

func Scrape(url string, depth int) ([]Page, map[string]string) {
	var scraped sync.Map
	var wg sync.WaitGroup
	var res []Page

	_, err := http.Get(url)
	if err != nil {
		return nil, map[string]string{"Error": "invalid url."}
	}

	scrape(url, depth, &scraped, &wg, &res)
	wg.Wait()

	return res, nil 
}
