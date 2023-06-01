package getreadmeurl

import "github.com/gocolly/colly"

type Readme interface {
	GetReadmeURL()
}

type ReadmeElement struct {
	TopPage string
	Target  string
}

var TargetURL []string

func SetReadmeURL(r Readme) {
	r.GetReadmeURL()
}

func (r ReadmeElement) GetReadmeURL() {

	c := colly.NewCollector()
	c.OnHTML(r.Target, func(h *colly.HTMLElement) {
		TargetURL = append(TargetURL, h.Attr("href"))
	})
	c.Visit(r.TopPage)

}