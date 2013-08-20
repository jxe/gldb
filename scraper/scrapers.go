package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jxe/gldb"
	"regexp"
	"strings"
)

func parentAttr(s *goquery.Selection, attr string) string {
	sel := "[" + attr + "]"
	res, _ := s.Closest(sel).Attr(sel)
	return res
}

func childElemValue(s *goquery.Selection, sel string) string {
	return s.ChildrenFiltered(sel).Text()
}

func htmlScraper(f *Fetcher) (recognized bool) {
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


func videoScraper(f *Fetcher) (recognized bool) {
	matched, _ := regexp.MatchString(f.url, "vimeo|youtube")
	if !matched {
		return false
	}
	f.fetched(&gldb.Doable{
		URL:         f.url,
		Type:        "video",
		Title:       "Unknown Video",
		Description: "",
		Metro:       "",
		Qualities:   []string{},
		Subjects:    []string{},
		GuideURLs:   []string{},
	})
	return true
}

func defaultScraper(f *Fetcher) (recognized bool) {
	f.fetched(&gldb.Doable{
		URL:         f.url,
		Type:        "article",
		Title:       "Unknown Article",
		Description: "",
		Metro:       "",
		Qualities:   []string{},
		Subjects:    []string{},
		GuideURLs:   []string{},
	})
	return true
}

var scrapers = []scraper{
	htmlScraper,
	videoScraper,
	defaultScraper,
}
