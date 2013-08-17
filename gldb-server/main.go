package main

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/jxe/gldb"
	"github.com/jxe/gldb/scraper"
	"log"
	"os"
	"strings"
)

func main() {
	log.Print("Starting GLDB server")
	var db, err = gldb.GLDBFromMongoURL(os.Getenv("GLDB_MONGO_URL"))
	if err != nil {
		panic(err)
	}
	log.Print("Connected to mongolab")


	// data collection

	web.Post("/did", func(c *web.Context) string {
		err := db.AddReview(&gldb.Review{
			DoableURL:            c.Params["what"],
			City:                 c.Params["city"],
			Comment:              c.Params["comment"],
			SatisfiedVibes:       strings.Split(c.Params["tws"], ","),
			AuthorURLs:           []string{},
			RelativeToDoableURLs: []string{},
		})
		if err != nil {
			panic(err)
		}
		return "thanks"
	})


	// querying

	web.Get("/skills", func(c *web.Context) string {
		skills := db.Skills(c.Params["city"], c.Params["vibe"])
		body, err := json.Marshal(skills)
		if err != nil {
			panic(err)
		}
		return string(body)
	})


	// debugging

	web.Get("/raw", func(c *web.Context) string {
		f := scraper.Fetcher{ url: x.Params["url"] }
		return string(f.ResponseBodyBytes())
	})

	web.Get("/doableJSON", func(c *web.Context) string {
		l := &gldb.Doable{}
		results := []scraper.Scrapable{}
		err := Slurp(c.Params["url"], l, *results)
		body, err := json.Marshal(l)
		return string(body)
	})

	web.Get("/guideJSON", func(c *web.Context) string {
		l := &gldb.Guide{}
		results := []scraper.Scrapable{}
		err := Slurp(c.Params["url"], l, *results)
		body, err := json.Marshal(l)
		return string(body)
	})


	web.Run("0.0.0.0:9999")
}
