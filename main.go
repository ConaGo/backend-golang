package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"conago.de/web-scraper/recurring_tasks"
	"conago.de/web-scraper/server"
)

//"conago.de/web-scraper/html_parser"
//"conago.de/web-scraper/data_parser"
//"conago.de/web-scraper/server"
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)
func main() {
	//html_parser.Test()
	//data_parser.ParseData()
	if len(os.Args) >1 &&os.Args[1] == "first"{
		recurring_tasks.GetConferenceData()
	}

	go doTask()
	server.Start()



}
func doTask() {
	time.Sleep(60*time.Second)
	recurring_tasks.GetConferenceData()
	os.Exit(0)
}