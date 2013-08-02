package main

import (
	"github.com/hoisie/web"
	"github.com/jxe/gldb"
    "encoding/json"
    "log"
    "strings"
)

func main() {
    log.Print("Starting GLDB server")
    var db, err = gldb.GLDBFromMongoURL("mongodb://gldb:popple@ds027668.mongolab.com:27668/gldb")
    if err != nil {
        panic(err)
    }
    log.Print("Connected to mongolab")

	web.Get("/did", func(c *web.Context) string {
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
