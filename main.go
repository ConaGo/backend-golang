package main

import (
	"path/filepath"
	"runtime"

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
	server.StartServer()
	//recurring_tasks.GetConferenceData()
}