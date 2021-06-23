package recurring_tasks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Init() {

}
func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
func executeCommand(command,  directory string, arguments ...string) {
	cmd := exec.Command(command, arguments...)
	if directory != "" {
		cmd.Dir = directory
	}
    out, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Printf("%s", out);
    }
}

func GetConferenceData() {
	//pull latest update from repo
	executeCommand("git", "data/conference-data", "pull")
	//create a temporary database to and parse the data
	newDB, err := os.Create("data/tmp.db")
	check(err)
	fmt.Println(newDB)
	//data_parser.ParseData("data/tmp.db")
	
	//swap out the temp db with the actual db
	//server.StopServer()
	executeCommand("rm", "data", "conference.db",)
	executeCommand("mv", "data", "tmp.db", "conference.db" )
	//server.StartServer()

}