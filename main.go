package main

import (
	"fmt"
	"net/http"

	"conago.de/myutil"
	"conago.de/web-scraper/server"

	//"conago.de/web-scraper/html_parser"
	//"conago.de/web-scraper/data_parser"
	//"conago.de/web-scraper/server"

	"github.com/gocolly/colly/v2"
)

func main() {
	//html_parser.Test()
	//data_parser.ParseData()
	http.HandleFunc("/", server.Serve)
	http.ListenAndServe(":8080", nil)

}

func secondary() {
	c := colly.NewCollector(
	//colly.AllowedDomains("https://confs.tech/", "www.confs.tech/", "https://www.confs.tech/"),
	//colly.AllowedDomains("about.gitlab.com", "https://about.gitlab.com/"),
	)

	//add print callback to scraper
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		//change user-agent field on every request
		//r.Headers.Set("User-Agent","Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4")
		r.Headers.Set("User-Agent", myutil.RandomString())
	})

	c.OnHTML("*", func(e *colly.HTMLElement) {
		//goquerySelection := e.DOM
		fmt.Println(e.Text)
		// Example Goquery usage
		//fmt.Println(goquerySelection.Find("div").Attr("class"))
	})
}
