package awesome

import (
	"collections/getreadmeurl"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

const TopPage = "https://github.com/matiassingers/awesome-readme/blob/master/readme.md"
const Target = "article.markdown-body.entry-content.container-lg > ul > li > a"

func InputToSection() []string {
	c := colly.NewCollector()
	const targetURL = "article.markdown-body.entry-content.container-lg >h2"
	s := []string{}
	c.OnHTML(targetURL, func(h *colly.HTMLElement) {
		s = append(s, h.DOM.Text())
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL:", r.URL.String())
	})
	for _, v := range getreadmeurl.TargetURL {
		c.Visit(v)
		time.Sleep(1 * time.Second)
	}
	return s
}

func CountSection(s []string) {

	total := s
	fmt.Printf("集計ReadMe数: %v \n", len(total))
	fmt.Println("---使用率---")
	counts := make(map[string]int)
	for _, v := range s {
		counts[v]++
	}
	for k, v := range counts {
		m := float32(len(total))
		n := float32(100 * v)
		q := n / m
		fmt.Printf("%s: %f%% %v/%v \n", k, q, v, len(total))
	}
}