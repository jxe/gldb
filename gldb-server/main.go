package main

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/jxe/gldb"
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

	web.Post("/did", func(c *web.Context) string {
		r := &gldb.Review{
			DoableURL:            c.Params["what"],
			City:                 c.Params["city"],
			Comment:              c.Params["comment"],
			SatisfiedDesires:     strings.Split(c.Params["tws"], ","),
			AuthorURLs:           []string{},
			RelativeToDoableURLs: []string{},
		}
		err := db.AddReview(r)
		if err != nil {
			panic(err)
		}
		return "thanks"
	})

	web.Get("/topics", func(c *web.Context) string {
		topics := db.Topics(c.Params["city"], c.Params["desire"])
		body, err := json.Marshal(topics)
		if err != nil {
			panic(err)
		}
		return string(body)
	})

	web.Run("0.0.0.0:9999")
}
