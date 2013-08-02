package gldb

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)



type GLDB struct {
	topics  *mgo.Collection
	doables *mgo.Collection
	guides  *mgo.Collection
	reviews *mgo.Collection
}




// TODO

func (t *TopicInCity) Doables() []Doable { return nil }


func doableFromURL(url string) (d *Doable, err error) {
	d = &Doable{}
	err = jsonFromURL(url, d)
	if err == nil {
		d.URL = url
	}
	return
}

func (d *Doable) loadGuides() (guides []*Guide, err error) {
	guides = make([]*Guide, len(d.GuideURLs))
	for i, url := range d.GuideURLs {
		g := &Guide{}
		err = jsonFromURL(url, g)
		if err != nil {
			return
		}
		g.URL = url
		guides[i] = g
	}
	return
}


func (db *GLDB) AddReview(r *Review) (err error) {
	doable, err := doableFromURL(r.DoableURL)
	if err != nil {
		return
	}
	guides, err := doable.loadGuides()
	if err != nil {
		return
	}
	db.ensureDoable(doable)
	for _, g := range guides {
		db.ensureGuide(g)
	}
	db.reviews.Insert(r)
	return
}



// MongoDB specific


func (db *GLDB) Topics(city, desire string) (result []*TopicInCity) {
	// for now, just return topics that have the city and desire
	// later, add pools of doables under a topic without a guide
	topics := db.TopicIDsForCityAndDesire(city, desire)
	result = make([]*TopicInCity, len(topics))
	for i, t := range topics {
		guides := []*Guide{}
		result[i] = &TopicInCity{city, t, guides}
		err := db.guides.Find(bson.M{"topic": t, "city": city}).All(&guides)
		if err != nil {
			panic("Failed to get guides")
		}
		result[i].Guides = guides
	}
	return
}


func (db *GLDB) TopicIDsForCityAndDesire(city, desire string) (topics []string) {
	err := db.guides.Find(bson.M{"desires": desire, "city": city}).Distinct("topic", &topics)
	if err != nil {
		panic("Failed to get topics")
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
		topics:  db.C("topics"),
		doables: db.C("doables"),
		guides:  db.C("guides"),
		reviews: db.C("reviews"),
	}
	return
}

func (db *GLDB) Close() {
	db.topics.Database.Session.Close()
	return 
}

func (db *GLDB) ensureDoable(d *Doable) {
	_, err := db.doables.Upsert(bson.M{"URL": d.URL}, d)
	if err != nil {
		panic(err)
	}
}

func (db *GLDB) ensureGuide(g *Guide) {
	_, err := db.guides.Upsert(bson.M{ "URL":g.URL }, g)
	if err != nil {
		panic(err)
	}
}
