package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"conago.de/myutil"
	"conago.de/web-scraper/data_parser"
	"conago.de/web-scraper/db/sqlite"
)

func sanitizeTags(values []string) error {
	allowedTags := []string{"javascript", "css", "devops", "leadership", "android", "clojure", "cpp", "data", "dotnet", "elixir", "general", "golang", "graphql", "ios", "iot", "java", "kotlin", "networking", "php", "product", "python", "ruby", "rust", "security", "tech-comm", "typescript", "ux", "cpp", "elixir", "elm", "groovy", "scala"}
	for _, value := range values {
		if !myutil.StringInSlice(value, allowedTags) {
			return errors.New("malformed tags")
		}
	}
	return nil
}
func HandleConferences(w http.ResponseWriter, r *http.Request) {
	//CORS - remove
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := r.URL.Query()
	tags := params["tag"]
	date := params.Get("date")
	t, dateErr := time.Parse(layoutISO, date)
	dateTime := time.Now()
	if dateErr == nil {
		dateTime = t
	}
	if tagErr := sanitizeTags(tags); tagErr != nil {
		log.Println(tagErr)
		return 
	}
	conferences := []data_parser.Conference{}
	_conferences := []data_parser.Conference{}
	if len(tags) == 0 {
		sqlite.ConfDB.DB.Preload("Tags").Where("end_date > ?", dateTime).Order("start_date ASC").Find(&_conferences)
	} else {
		sqlite.ConfDB.DB.Preload("Tags", "tag_name IN ?", tags).Where("end_date > ?", dateTime).Order("start_date ASC").Find(&_conferences)
	}
	for _, c := range _conferences {
		if len(c.Tags) > 0 {
			conferences = append(conferences, c)
		}
	}
	if err := sqlite.ConfDB.DB.Error; err != nil {
        log.Println(err)
		return
    }
	b, err := json.Marshal(conferences); if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}