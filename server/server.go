package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"conago.de/web-scraper/data_parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test() {
	//https://golang.org/pkg/time/#example_Tick
	t := time.Tick(30 * time.Minute)
	for next := range t {
		fmt.Printf("%v %s\n", next, statusUpdate())
	}
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//tags := params["tags"]
	tag := params.Get("tag")
	fmt.Println(tag)
	db, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var conferences []data_parser.Conference
	db.Raw("SELECT * FROM tags INNER JOIN conference_tags ON tags.tag_id = conference_tags.tag_tag_id AND tags.tag_name = ? INNER JOIN conferences ON conferences.id = conference_tags.conference_id", tag).Scan(&conferences)

	b, err := json.Marshal(conferences)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}
