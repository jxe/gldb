package gldb

type Doable struct {
	URL           string
	Kind          string
	Title         string
	City          string
	Desires       []string
	Topics        []string
	GuideURLs     []string
	AuthorURLs    []string
	ProvidersURLs []string
}

type Guide struct {
	URL        string
	Title      string
	City       string
	Topic      string
	AuthorURLs []string
	Desires    []string
}

type Review struct {
	DoableURL            string
	City                 string
	Comment              string
	AuthorURLs           []string
	SatisfiedDesires     []string
	RelativeToDoableURLs []string
}

type Topic struct {
	id         string
	aliasOf    string
	FullName   string
	otherNames []string
}

type TopicInCity struct {
	City   string
	Topic  string
	Guides []*Guide
}
