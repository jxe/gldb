package gldb

import "time"

type GLDB interface {
	AddReviewAndRelatedData(r *Review, metro, sociographic, comment string, loadables []interface{})
	SubjectsInMetro(metro string, quality string) (result []*SubjectInMetro)
}

type Doable struct {
	URL           string
	Type          string
	Title         string
	Description   string
	Metro         string
	Qualities     []string
	Subjects      []string
	GuideURLs     []string
	ProvidersURLs []string
}

type Guide struct {
	URL         string
	Title       string
	Description string
	Metro       string
	Subject     string
	Qualities   []string
}

type ReviewPeriod struct {
	ReviewerURL           string
	ReviewPeriodStartTime time.Time
	ReviewPeriodEndTime   time.Time
	Metro                 string
	SociographicData      string
	Comment               string
}

type Review struct {
	DoableURL             string
	ReviewerURL           string
	ReviewPeriodStartTime time.Time
	ReviewPeriodEndTime   time.Time
	ReviewTime            time.Time
	QualitiesConfirmed    []string
	Comment               string
}

type Subject struct {
	id         string
	aliasOf    string
	FullName   string
	otherNames []string
}

type SubjectInMetro struct {
	Metro   string
	Subject string
	Guides  []*Guide
}
