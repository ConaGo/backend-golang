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
	http.HandleFunc("/conferences", HandleConferences)
	http.HandleFunc("/questions", HandleQuestions)
	http.ListenAndServe(":8080", nil)
}

