package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"conago.de/myutil"
	"conago.de/web-scraper/data_parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
const (
	layoutISO = "2006-01-02"
)
func Serve() {
/* 	//https://golang.org/pkg/time/#example_Tick
	t := time.Tick(30 * time.Minute)
	for next := range t {
		fmt.Printf("%v %s\n", next, statusUpdate())
	} */
	http.HandleFunc("/conferences", handleConferences)
	http.HandleFunc("/questions")
	http.ListenAndServe(":8080", nil)
}
func sanitizeParams(name, value string) error{
	allowedTags := []string{"javascript", "css", "devops", "leadership", "android", "clojure", "cpp", "data", "dotnet", "elixir", "general", "golang", "graphql", "ios", "iot", "java", "kotlin", "networking", "php", "product", "python", "ruby", "rust", "security", "tech-comm", "typescript", "ux", "cpp", "elixir", "elm", "groovy", "scala"}
	switch name {
	case "tag":
		if myutil.StringInSlice(value, allowedTags){
			return nil
		}
	}
	return errors.New("malformed params")
}
func handleConferences(w http.ResponseWriter, r *http.Request) {
	//CORS - remove 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	params := r.URL.Query()
	tags := params["tag"]
	date := params.Get("date")
	fmt.Println(date)
	t, err := time.Parse(layoutISO, date)
	dateTime := time.Now()
	if err == nil {
		dateTime = t
	}
	db, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	conferences := []data_parser.Conference{}
	if len(tags) == 0 {
		db.Debug().Preload("Tags").Where("end_date > ?", dateTime).Order("start_date ASC").Find(&conferences)
	} else {
		db.Debug().Preload("Tags", "tag_name IN ?", tags).Where("end_date > ?", dateTime).Order("start_date ASC").Find(&conferences)
	}
		

	var _conferences []data_parser.Conference
	for _, c:= range conferences {
		if len(c.Tags) > 0 {
			_conferences = append(_conferences, c)
		} 
	}
	b, err := json.Marshal(_conferences)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}
func 