package gldb

import (
	"github.com/jxe/gldb/scraper"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type GLDB struct {
	skills  *mgo.Collection
	doables *mgo.Collection
	guides  *mgo.Collection
	reviews *mgo.Collection
}

// TODO

type loadable interface {
	loadInDB(*GLDB)
}

func (db *GLDB) AddReview(r *Review) (err error) {
	scrapeResults, err := scraper.Slurp(r.DoableURL)
	if err != nil {
		return
	}
	for _, r := range scrapeResults {
		r.(loadable).loadInDB(db)
	}
	return
}

type SkillInCity struct {
	City   string
	Skill  string
	Guides []*Guide
}

func (t *SkillInCity) Doables() []Doable { return nil }

// MongoDB specific

func (db *GLDB) Skills(city, vibe string) (result []*SkillInCity) {
	// for now, just return skills that have the city and vibe
	// later, add pools of doables under a skill without a guide
	skills := db.SkillIDsForCityAndVibe(city, vibe)
	result = make([]*SkillInCity, len(skills))
	for i, t := range skills {
		guides := []*Guide{}
		result[i] = &SkillInCity{city, t, guides}
		err := db.guides.Find(bson.M{"skill": t, "city": city}).All(&guides)
		if err != nil {
			panic("Failed to get guides: " + err.Error())
		}
		result[i].Guides = guides
	}
	return
}

func (db *GLDB) SkillIDsForCityAndVibe(city, vibe string) (skills []string) {
	err := db.guides.Find(bson.M{"vibes": vibe, "city": city}).Distinct("skill", &skills)
	if err != nil {
		panic("Failed to get skills" + err.Error())
	}
	return
}

func GLDBFromMongoURL(url string) (d *GLDB, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return
	}
	db := session.DB("")
	d = &GLDB{
		skills:  db.C("skills"),
		doables: db.C("doables"),
		guides:  db.C("guides"),
		reviews: db.C("reviews"),
	}
	return
}

func (db *GLDB) Close() {
	db.skills.Database.Session.Close()
	return
}

func (r *Review) loadInDB(db *GLDB) {
	db.reviews.Insert(r)
}

func (g *Guide) loadInDB(db *GLDB) {
	_, err := db.guides.Upsert(bson.M{"URL": g.URL}, g)
	if err != nil {
		panic(err)
	}
}

func (d *Doable) loadInDB(db *GLDB) {
	_, err := db.guides.Upsert(bson.M{"URL": d.URL}, d)
	if err != nil {
		panic(err)
	}
}
