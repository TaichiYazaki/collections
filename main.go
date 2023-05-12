package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gocolly/colly"
)

type Collection struct {
	name  string  
	value int	
}

type List []Collection

func getReadmeURL() []string {
	c := colly.NewCollector()
	u := "article.markdown-body.entry-content.container-lg > ul > li > a"
	var url []string
	c.OnHTML(u, func(h *colly.HTMLElement) {
		url = append(url, h.Attr("href"))
	})
	c.Visit("https://github.com/matiassingers/awesome-readme/blob/master/readme.md")

	return url
}

func sectionToCsv() {
	c := colly.NewCollector()
	u := "article.markdown-body.entry-content.container-lg >h2"
	if _, err := os.Stat("section.csv"); err == nil {
		r := os.Remove("section.csv")
		fmt.Println(r, "section.csvをリセット")
	}
	c.OnHTML(u, func(h *colly.HTMLElement) {
		var section []string
		sections := append(section, h.DOM.Text())
		f, err := os.OpenFile("section.csv", os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w := csv.NewWriter(f)
		w.Write(sections)
		w.Flush()
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL:", r.URL.String())
	})
	for _, v := range getReadmeURL() {
		c.Visit(v)
		time.Sleep(1 * time.Second)
	}
}

func sortSection() {
	f, err := os.Open("section.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	collections, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	counts := make(map[string]int)
	for _, c := range collections {
		for _, v := range c {
			counts[v]++
		}
	}
	l :=List{}
	for k, v := range counts {
		c := Collection{k, v}
		l = append(l,c)
	}
	sort.Sort(l)
	for _, v := range l {
		fmt.Println(v)
	}
}

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	} else {
		return (l[i].value > l[j].value)
	}
}

func main() {
	sectionToCsv()
	sortSection()
}
