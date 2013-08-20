package scraper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Slurp(u string) (results []interface{}, err error) {
	defer func() {
		thing := recover()
		switch thing.(type) {
		case error:
			err = thing.(error)
		}
	}()
	f := &Fetcher{urlsToFetch: []string{u}}
	f.spin()
	return f.fetchedObjects, err
}

type scraper func(f *Fetcher) (recognized bool)

type Fetcher struct {
	url            string
	body           *[]byte
	contentType    *string
	response       *http.Response
	html           *goquery.Document
	urlsToFetch    []string
	fetchedObjects []interface{}
}

func (f *Fetcher) fetched(thing interface{}) {
	f.fetchedObjects = append(f.fetchedObjects, thing)
}

func (f *Fetcher) spin() {
FETCHES:
	for len(f.urlsToFetch) > 0 {
		f.url = f.urlsToFetch[0]
		f.urlsToFetch = f.urlsToFetch[1:]
		f.body, f.contentType, f.response, f.html = nil, nil, nil, nil

		for _, scraper := range scrapers {
			if scraper(f) {
				continue FETCHES
			}
		}
		panic(errors.New("No recognizer found"))
		return
	}
	return
}

func NewFetcherForURL(url string) *Fetcher {
	return &Fetcher{url: url}
}

func (f *Fetcher) Response() (r *http.Response) {
	if f.response == nil {
		var err error
		f.response, err = http.Get(f.url)
		if err != nil {
			panic(err)
		}
	}
	return f.response
}

func (f *Fetcher) ResponseBodyBytes() (body *[]byte) {
	if f.body == nil {
		defer f.Response().Body.Close()
		var err error
		body, err := ioutil.ReadAll(f.Response().Body)
		if err != nil {
			panic(err)
		}
		f.body = &body
	}
	return f.body
}

func (f *Fetcher) ContentType() string {
	return f.Response().Header.Get("Content-Type")
}

func (f *Fetcher) HTML() *goquery.Document {
	if f.html == nil {
		html, err := goquery.NewDocumentFromResponse(f.Response())
		if err != nil {
			panic(err)
		}
		f.html = html
	}
	return f.html
}

func (f *Fetcher) ParsedURL() *url.URL {
	parsed, err := url.Parse(f.url)
	if err != nil {
		panic(err)
	}
	return parsed
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
