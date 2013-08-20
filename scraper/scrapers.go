package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jxe/gldb"
	"strings"
)

// function findDoables(){
//   // gather anything with a [doable] attr
//   // search for subject, quality, and metro in self/parents
//   // search for title and description as child headers and paras
//   // check if guides in parent
// }

// function findGuides(){
//   // gather anything with a [subject] that's not [doable]
// }

func parentAttr(s *goquery.Selection, attr string) string {
	return ""
}

func childElemValue(s *goquery.Selection, sel string) string {
	return ""
}

func htmlScraper(f Fetcher) (recognized bool) {
	e := f.HTML()
	foundSomething := false

	e.Find("[doable]").Each(func(i int, doable *goquery.Selection) {
		foundSomething = true
		t, _ := doable.Attr("doable")
		id, _ := doable.Attr("id")
		f.fetched(&gldb.Doable{
			URL:         f.url + "#" + id,
			Type:        t,
			Title:       childElemValue(doable, "h3"),
			Description: childElemValue(doable, "p"),
			Metro:       parentAttr(doable, "metro"),
			Qualities:   strings.Split(parentAttr(doable, "qualities"), ","),
			Subjects:    strings.Split(parentAttr(doable, "subjects"), ","),
			GuideURLs:   []string{f.url},
		})
	})

	e.Find("[subject]").Not("[doable]").Each(func(i int, guide *goquery.Selection) {
		foundSomething = true
		s, _ := guide.Attr("subjects")
		f.fetched(&gldb.Guide{
			URL:         f.url,
			Subject:     s,
			Title:       childElemValue(guide, "h3"),
			Description: childElemValue(guide, "p"),
			Metro:       parentAttr(guide, "metro"),
			Qualities:   strings.Split(parentAttr(guide, "qualities"), ","),
		})
	})

	return foundSomething
}

var scrapers = []scraper{
	htmlScraper,
	// jsonScraper,
}
