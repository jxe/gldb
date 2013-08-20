package gldb

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoGLDB struct{ *mgo.Database }

func (db MongoGLDB) AddReviewAndRelatedData(r *Review, metro string, sociographic string, comment string, loadables []interface{}) {
	for _, obj := range loadables {
		switch obj.(type) {
		case *Guide:
			db.C("guides").Upsert(bson.M{"URL": obj.(*Guide).URL}, obj)
		case *Doable:
			db.C("doables").Upsert(bson.M{"URL": obj.(*Doable).URL}, obj)
		}
	}
	db.C("reviewPeriods").Upsert(bson.M{
		"ReviewerURL":           r.ReviewerURL,
		"ReviewPeriodStartTime": r.ReviewPeriodStartTime,
		"ReviewPeriodEndTime":   r.ReviewPeriodEndTime,
	}, &ReviewPeriod{
		ReviewerURL:           r.ReviewerURL,
		ReviewPeriodStartTime: r.ReviewPeriodStartTime,
		ReviewPeriodEndTime:   r.ReviewPeriodEndTime,
		Metro:                 metro,
		SociographicData:      sociographic,
		Comment:               comment,
	})
	db.C("reviews").Insert(r)
	return
}

func (db MongoGLDB) SubjectsInMetro(metro, quality string) (result []*SubjectInMetro) {
	subjects := []string{}
	err := db.C("guides").Find(bson.M{"qualities": quality, "metro": metro}).Distinct("subject", &subjects)
	if err != nil {
		panic("Failed to get skills" + err.Error())
	}
	result = make([]*SubjectInMetro, len(subjects))
	for i, t := range subjects {
		guides := []*Guide{}
		result[i] = &SubjectInMetro{metro, t, guides}
		err := db.C("guides").Find(bson.M{"subject": t, "metro": metro}).All(&guides)
		if err != nil {
			panic("Failed to get guides: " + err.Error())
		}
		result[i].Guides = guides
	}
	return
}

func (db MongoGLDB) DoablesForSubjectInMetro(metro, subject string) (result []*Doable) {
	err := db.C("doables").Find(bson.M{"metro": metro, "subjects": subject}).All(&result)
	if err != nil {
		panic(err)
	}
	return
}


func GLDBFromMongoURL(url string) (d GLDB, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return
	}
	return MongoGLDB{session.DB("")}, nil
}

func (db MongoGLDB) Close() {
	db.Session.Close()
	return
}
