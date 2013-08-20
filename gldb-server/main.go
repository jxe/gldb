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
		r := &gldb.Review{
			DoableURL:          c.Params["what"],
			Comment:            c.Params["comment"],
			QualitiesConfirmed: strings.Split(c.Params["qualities"], ","),
		}

		scrapeResults, err := scraper.Slurp(r.DoableURL)
		if err != nil {
			panic(err)
		}

		db.AddReviewAndRelatedData(r, c.Params["metro"], c.Params["sociographic"], c.Params["comment"], scrapeResults)
		return "thanks"
	})

	// querying

	web.Get("/subjects", func(c *web.Context) string {
		c.ContentType("json")
		interests := db.SubjectsInMetro(c.Params["metro"], c.Params["quality"])
		body, err := json.MarshalIndent(interests, "", "    ")
		if err != nil {
			panic(err)
		}
		return string(body)
	})

	web.Get("/subject", func(c *web.Context) string {
		c.ContentType("json")
		doables := db.DoablesForSubjectInMetro(c.Params["metro"], c.Params["subject"])
		body, err := json.MarshalIndent(doables, "", "    ")
		if err != nil {
			panic(err)
		}
		return string(body)
	})

	// debugging

	web.Get("/raw", func(c *web.Context) string {
		f := scraper.NewFetcherForURL(c.Params["url"])
		return string(*f.ResponseBodyBytes())
	})

	web.Get("/json", func(c *web.Context) string {
		c.ContentType("json")
		results, err := scraper.Slurp(c.Params["url"])
		if err != nil {
			panic(err)
		}
		body, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			panic(err)
		}
		return string(body)
	})

	web.Run("0.0.0.0:9999")
}
