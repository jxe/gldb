package scraper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Slurp(u string) (results []interface{}, err error) {
	defer func() {
		err = recover()
	}()
	f := &Fetcher{urlsToFetch: []string{u}}
	f.spin()
	return f.fetchedObjects, err
}

type scraper func(f Fetcher) (recognized bool)

type Fetcher struct {
	url            string
	body           *[]byte
	contentType    *string
	response       *http.Response
	html           *goquery.Document
	urlsToFetch    []string
	fetchedObjects []interface{}
}

func (f Fetcher) spin() {
FETCHES:
	for len(f.urlsToFetch) > 0 {
		f.url = urlsToFetch[0]
		f.urlsToFetch = f.urlsToFetch[1:]
		f.body, f.contentType, f.response, f.html = nil, nil, nil, nil

		for _, scraper := range scrapers {
			if scraper(f) {
				continue FETCHES
			}
		}
		panic("No recognizer found")
		return
	}
	return
}

func (f Fetcher) Response() (r *http.Response) {
	if !f.response {
		f.response, err = http.Get(f.url)
		if err != nil {
			panic(err)
		}
	}
	return f.response
}

func (f Fetcher) ResponseBodyBytes() (body []byte) {
	if !f.body {
		defer f.Response().Body.Close()
		f.body, err = ioutil.ReadAll(f.Response().Body)
		if err != nil {
			panic(err)
		}
	}
	return f.body
}

func (f Fetcher) ContentType() string {
	f.Response().Header.Get("Content-Type")
}

func (f Fetcher) HTML() *goquery.Document {
	if !f.html {
		f.html, err = goquery.NewDocumentFromResponse(f.Response())
	}
	return f.html
}

func (f Fetcher) ParsedURL() *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return parsed
}

func microFormatScraper(f Fetcher) (recognized bool) {
	// freak out if it's not html
	e := f.HTML()
	g := e.Find(".guide")
	switch l.(type) {
	case *Doable:
		fragment := f.ParsedURL().Fragment
		if len(fragment) > 0 {
			e = doc.Find("#" + fragment)
		} else {
			e = doc.Find("[data-doable]")
		}

		l.URL = u
		l.Kind = e.Attr("data-doable")
		l.City = e.Attr("data-region")
		l.Vibes = strings.Split(e.Attr("data-vibe", ","))
		l.Skills = strings.Split(e.Attr("data-skills", ";"))
		l.Title = e.Find(".doable-title").Text()
		l.Notes = e.Find(".doable-notes").Text()
		l.GuideURLs = e.Find("a[rel=guide]").Map(func(_, s *goquery.Selection) string {
			return s.Attr("href")
		})

		// configure guide defaults
		if g {
			if len(l.City) == 0 {
				l.City = g.Attr("data-region")
			}
			if len(l.Vibes) == 0 {
				l.Vibes = strings.Split(g.Attr("data-vibe", ","))
			}
			if len(l.Skills) == 0 {
				l.Skills = strings.Split(g.Attr("data-skills"), "; ")
			}
			if len(l.GuideURLs) == 0 {
				l.GuideURLs = []string{strings.Split(u, '#')[0]}
			}
		}

		// if there's just one h3, that's the title
		if len(l.Title) == 0 && e.Find("h3") {
			l.Title = e.Find("h3").Text()
		}

		// if there's just one p, that's the notes
		if len(l.Notes) == 0 && e.Find("p") {
			l.Title = e.Find("p").Text()
		}

		return true, nil

	case *Guide:
		l.URL = u
		l.City = g.Attr("data-region")
		l.Vibes = strings.Split(g.Attr("data-vibe", ","))
		l.Title = g.Find(".guide-title").Text()
		l.Skill = strings.Split(g.Attr("data-skills"), "; ")

		if len(l.Title) == 0 && g.Find("h1") {
			l.Title = g.Find("h1").Text()
		}

		return true, nil
	}
}

var scrapers = []scraper{
	microFormatScraper,
	// jsonScraper,
}

// func jsonScraper(u string, l Scrapable, f Fetcher) (recognized bool, err error) {
// 	if !strings.Contains(u, ".json") || f.ContentType != "application/json" {
// 		return nil, nil
// 	}

// 	body := f.ResponseBodyBytes()
// 	if body[0] != '{' {
// 		return nil, nil
// 	}

// 	fragment := f.ParsedURL().Fragment
// 	if len(fragment) > 0 {
// 		var f interface{}
// 		err = json.Unmarshal(body, &f)
// 		if err != nil {
// 			return
// 		}
// 		m := f.(map[string]interface{})
// 		body, err = json.Marshal(m[fragment])
// 		if err != nil {
// 			return
// 		}
// 	}
// 	err = json.Unmarshal(body, l)
// 	return

// }
