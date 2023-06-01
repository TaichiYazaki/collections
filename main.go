package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gocolly/colly"
)

type Collection struct {
	Name  string  
	Value int	
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
	// CSVをリセットするためのコード開始
	if _, err := os.Stat("section.csv"); err == nil {
		r := os.Remove("section.csv")
		fmt.Println(r, "section.csvをリセット")
	}
	// 終了
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

func sortSectionToCsv() {
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
	
	var c []Collection
	for _, v := range l {
		l := Collection(v)
		c = append(c, l)
	}
	// CSVをリセットするためのコード開始
	if _, err := os.Stat("result.csv"); err == nil {
		r := os.Remove("result.csv")
		fmt.Println(r, "result.csvをリセット")
	}
	// 終了
	fi, err := os.OpenFile("reult.csv", os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err !=nil{
		log.Fatal(err)
	}
	defer fi.Close()
	gocsv.MarshalFile(c, fi)
}

//並べ替えるコレクションの要素数を取得
func (l List) Len() int {
	return len(l)
}

// 隣り合うに要素でどのような条件を満たした場合にSwapメソッドを実行するのか定義
func (l List) Less(i, j int) bool {
	if l[i].Value == l[j].Value {
		return (l[i].Name < l[j].Name)
	} else {
		return (l[i].Value > l[j].Value)
	}
}
// 隣り合う二要素に対してソートで入れ替えを行う際にどのような入れ替えを行うのか定義
func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func main() {
	sectionToCsv()
	sortSectionToCsv()
}
