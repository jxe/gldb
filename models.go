package gldb

type Doable struct {
	URL           string
	Kind          string
	Title         string
	City          string
	Notes		  string
	Vibes         []string
	Skills        []string
	GuideURLs     []string
	AuthorURLs    []string
	ProvidersURLs []string
}

type Guide struct {
	URL        string
	Title      string
	City       string
	Skill      string
	AuthorURLs []string
	Vibes      []string
}

type Review struct {
	DoableURL            string
	City                 string
	Comment              string
	AuthorURLs           []string
	SatisfiedVibes       []string
	RelativeToDoableURLs []string
}

type Skill struct {
	id         string
	aliasOf    string
	FullName   string
	otherNames []string
}
