package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"conago.de/myutil"
	"conago.de/web-scraper/data_parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test() {
/* 	//https://golang.org/pkg/time/#example_Tick
	t := time.Tick(30 * time.Minute)
	for next := range t {
		fmt.Printf("%v %s\n", next, statusUpdate())
	} */
	http.HandleFunc("/", handleIndex)
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
func handleIndex(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//tags := params["tags"]
	tagParam := params.Get("tag")

	fmt.Println(tagParam)
	db, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var conferences []data_parser.Conference
	if tagParam != "" {
		//db.Raw("SELECT * FROM tags INNER JOIN conference_tags ON tags.tag_id = conference_tags.tag_tag_id AND tags.tag_name = ? INNER JOIN conferences ON conferences.id = conference_tags.conference_id", tag).Scan(&conferences)
	}else{
		//db.Raw("SELECT * FROM tags INNER JOIN conference_tags ON tags.tag_id = conference_tags.tag_tag_id INNER JOIN conferences ON conferences.id = conference_tags.conference_id").Scan(&conferences)
	}
	//db.Model(&data_parser.Conference{}).Joins("left join conference_tags on conferences.id = conference_tags.conference_id").Joins("JOIN tags ON tags.tag_id = conference_tags.tag_tag_id AND tags.tag_name = ?", "javascript").Scan(&conferences)
	//db.Model(&data_parser.Conference{}).Joins("JOIN tags ON tags.tag_id = conference_tags.tag_tag_id AND conferences.id = conference_tags.conference_id").Scan(&conferences)
	//db.Preload("Tags").Find(&conferences)
	tag := []data_parser.Tag{}
	db.Preload("Tags").Find(&conferences)
	//db.Where("tag_name IN ?", []string{"javascript", "css", "devops"}).Preload("Conferences").Find(&tag)
	fmt.Println(tag)
	b, err := json.Marshal(conferences)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}