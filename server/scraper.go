package main

import (
	"log"
	"net/http"

	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/DavidBelicza/TextRank/v2"
	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	Url 	string
	Title	string
	Summary []string
}

func getLinks(doc *goquery.Document, scraped *sync.Map, depth int, wg *sync.WaitGroup, result *[]Page) { //[]string {
	traversed := make(map[string]bool)
	doc.Find(".mw-parser-output a").Each(func(i int, s *goquery.Selection) {
		if len(traversed) >= 5 { return }

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

			_, _scraped := scraped.Load(link.String())
			if !_scraped && !traversed[link.String()] {
				traversed[link.String()] = true

				wg.Add(1)
				go func(l string) {
					scrape(l, depth-1, scraped, wg, result)
					wg.Done()
				}(link.String()[1:])
			}
		}
	})
}

var rule = textrank.NewDefaultRule()
var lang = textrank.NewDefaultLanguage()
var algorithm = textrank.NewDefaultAlgorithm()
// used to remove references and \n
var ref_reg = regexp.MustCompile(`(\[[0-9]+\])|(\n)`)
// gets the title from a given page url
var title_reg = regexp.MustCompile(`wiki/`)

func getSummary(doc *goquery.Document) []string {
	text := doc.Find(".mw-parser-output p").Text()

	tr := textrank.NewTextRank() 
	tr.Populate(text, lang, rule)
	sentences := textrank.FindSentencesByRelationWeight(tr, 2)

	var res []string
	for i, s := range sentences {
		if i >= 2 { break }	
		res = append(res, ref_reg.ReplaceAllString(s.Value, ""))
	}

	return res
}

func scrape(page string, depth int, scraped *sync.Map, wg *sync.WaitGroup, result *[]Page) {
	if depth == 0 { return }

	// check if we have already traversed this page
	_, contains := scraped.Load(page)
	if contains {
		return 
	}
	scraped.Store(page, true)
		
	// get the page html
	res, err := http.Get("https://en.m.wikipedia.org/" + page)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("Error: could not scrape https://en.m.wikipedia.org/wiki/" + page)
		return 
	}

	// load the html query document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var curr Page
	curr.Url = "https://en.m.wikipedia.org/" + page;
	// curr.Title = doc.Find(".mw-page-title-main").Text()
	curr.Title = strings.Title( title_reg.ReplaceAllLiteralString(page, "") )
	getLinks(doc, scraped, depth, wg, result)
	curr.Summary = getSummary(doc)

	*result = append(*result, curr)
}

func Scrape(url string, depth int) []Page {
	var scraped sync.Map
	var wg sync.WaitGroup
	var res []Page

	scrape(url, depth, &scraped, &wg, &res)
	wg.Wait()
	return res 
}

