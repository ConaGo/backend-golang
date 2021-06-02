package data_parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
type Conference struct {
	gorm.Model
	Name string
	Url string
	StartDate time.Time
	EndDate time.Time
	City string
	Country string
	Online bool
	Twitter string
	Tags [] string
}
func main2(){
	_, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
	  panic("failed to connect database")
	}
	//db.AutoMigrate(&Conference{})
	//readConferenceData()
	datas := convertConferenceData(readConferenceData())
	fmt.Println(datas)
	//db.Create(&datas)

/* 	var conferences []Conference
	db.Find(&conferences)
	for val := range conferences{
		fmt.Println(val)
	} */
}
// Converts parsed json from "map[string]interface" format to the "Conference" struct
func convertConferenceData(d []map[string]interface{}) []Conference{
	var confs []Conference
	const( layoutISO = "2006-01-02")
	for _, data := range d {
		c := Conference{
			Name: "",
			Url: "",
			StartDate: time.Now(),
			EndDate: time.Now(),
			City: "",
			Country: "",
			Online: false,
			Twitter: "",
			Tags: []string{},
		}
		for key := range data {
			switch key{
			case "name":
				c.Name = reflect.ValueOf(data["name"]).String()
			case "url":
				c.Url = reflect.ValueOf(data["url"]).String()
			case "startDate":
				str := reflect.ValueOf(data["startDate"]).String()
				c.StartDate, _ = time.Parse(layoutISO, str)
			case "endDate":
				str := reflect.ValueOf(data["endDate"]).String()
				c.EndDate, _ = time.Parse(layoutISO, str)
			case "city":
				c.City = reflect.ValueOf(data["city"]).String()
			case "country":
				c.Country = reflect.ValueOf(data["country"]).String()
			case "online":
				if reflect.ValueOf(data["online"]).Bool(){
					c.Online = true
				}else { c.Online = false}
			case "twitter":
				c.Twitter = reflect.ValueOf(data["twitter"]).String()
			case "tags":
				c.Tags = strings.Fields(reflect.ValueOf(data["twitter"]).String())
			}
			
		}
		confs = append(confs, c)
	}
	return confs
}
func readConferenceData() []map[string]interface{} {
	rootPath := "data/conference-data/conferences"
	var parsedJson []map[string]interface{};
	var currentDir string
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error{
		if err != nil {
            log.Fatalf(err.Error())
        }
		if info.IsDir(){ 
			currentDir = info.Name()
			fmt.Printf("directory name: %s\n", currentDir)
		} else if currentDir == "2022"{
			path := rootPath + "/" + currentDir + "/" + info.Name()
			fmt.Println(path)
			jsonFile, e := os.Open(path)
			if e != nil {
				fmt.Println(e)
			}
			// defer the closing of our jsonFile so that we can parse it later on
			defer jsonFile.Close()
			// read our opened xmlFile as a byte array.
			byteJson, _ := ioutil.ReadAll(jsonFile)

			var newParsedJson []map[string]interface{};
			json.Unmarshal([]byte(byteJson), &newParsedJson)

			// adding the filename as a tag : because the raw data encodes the tags through the filename
			for _, re := range newParsedJson {
				// remove .json file ending
				arr := strings.Split(info.Name(), ".")[0]
				// add new "tags" field to the elements
				re["tags"] = arr
			}
			parsedJson = append(parsedJson, newParsedJson...)
		}
		//fmt.Println(ress)
        return nil
    })

	// The data contains duplicate entries
	// (because tags are stored in the filenames in the original data), 
	// so we remove those while collecting the "tags"
	var parsedJsonAsSet []map[string]interface{}
	for i , val1 := range parsedJson {
		shouldAdd := true
		fmt.Println(reflect.ValueOf(val1["name"]).String())
		for j, val2 := range parsedJson{
			if j > i && compare(val1, val2) {
				val2["tags"] = reflect.ValueOf(val2["tags"]).String() + " " + reflect.ValueOf(val1["tags"]).String()
				fmt.Println(reflect.ValueOf(val2["tags"]).String())
				shouldAdd = false
			}
		}
		if shouldAdd {
			parsedJsonAsSet = append(parsedJsonAsSet, val1)
		}
	}
	for _, val := range parsedJsonAsSet {
		fmt.Println(reflect.ValueOf(val["tags"]).String())
	}
	return parsedJsonAsSet;
}
// the compare function returns true if "name" and "startDate" field is equal
func compare (val1, val2 map[string]interface{}) bool {
	return reflect.ValueOf(val1["name"]).String() == reflect.ValueOf(val2["name"]).String() && reflect.ValueOf(val1["startDate"]).String() == reflect.ValueOf(val2["startDate"]).String()
}