package data_parser

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"conago.de/web-scraper/html_parser"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Conference struct {
	gorm.Model	`json:"-"`
	Name      string               //`json:"name"`
	Url       string               //`json:"url"`
	StartDate JSONDate             //`json:"startDate"`
	EndDate   JSONDate             //`json:"endDate"`
	City      string               //`json:"city"`
	Country   string               //`json:"country"`
	Online    bool                 //`json:"online"`
	Twitter   string               //`json:"twitter"`
	Tags      []Tag                `gorm:"many2many:conference_tags;"`
	Metadata  html_parser.HTMLMeta `gorm:"embedded"`
}

func (me Conference) String() string {
	return fmt.Sprintf("Conference(name=%s,startDate=%v,tags=%v)\n", me.Name, me.StartDate, me.Tags)
}

const (
	layoutISO = "2006-01-02"
)

type JSONDate time.Time

func (j *JSONDate) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*j = JSONDate(t)
	}
	return nil
}

func (j JSONDate) Value() (driver.Value, error) {
	return time.Time(j), nil
}

// https://golang.org/pkg/encoding/json/#Unmarshaler
func (me *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(layoutISO, s)
	if err != nil {
		return err
	}
	*me = JSONDate(t)
	return nil
}
func (me JSONDate) MarshalJSON() ([]byte, error) {
	t := time.Time(me)
	stamp := fmt.Sprintf("\"%v-%v-%v\"",t.Year(), t.Month(), t.Day() )
    return []byte(stamp), nil
}

func (me JSONDate) String() string {
	return time.Time(me).String()
}

type Tag struct {
	gorm.Model				`json:"-"`
	TagName   	string
	Conferences []Conference `gorm:"many2many:conference_tags;"`
}

/* func (c *Conference) BeforeCreate(tx *gorm.DB) (err error) {
	c.UUID = uuid.New()
	return nil
  } */
func ParseData() {
	db, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Conference{})
	db.AutoMigrate((&Tag{}))
	datas := readConferenceData()
	db.Create(&datas)
	var conferences []Conference
	for _, elem := range conferences {
		fmt.Println(elem.Name, elem.City)
	}
}


func readConferenceData() []*Conference {
	rootPath := "data/conference-data/conferences"
	var allConferences []*Conference
	var currentDir string
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if info.IsDir() {
			currentDir = info.Name()
			fmt.Printf("directory name: %s\n", currentDir)
		} else if currentDir == "2022" {
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

			//* var newParsedJson []map[string]interface{}
			//* json.Unmarshal([]byte(byteJson), &newParsedJson)
			var fileConferences []*Conference
			json.Unmarshal(byteJson, &fileConferences)

			// adding the filename as a tag : because the raw data encodes the tags through the filename
			for _, conf := range fileConferences {
				// remove .json file ending
				tag := strings.Split(info.Name(), ".")[0]
				// add new "tags" field to the elements
				conf.Tags = []Tag{
					{
						TagName: tag,
					},
				}
			}
			allConferences = append(allConferences, fileConferences...)
		}
		return nil
	})

	// The data contains duplicate entries
	// (because tags are stored in the filenames in the original data),
	// so we remove those while collecting the "tags"
	var filteredConferences []*Conference
	for i, conference := range allConferences {
		shouldAdd := true
		for j, c := range allConferences {
			if j > i && compare(conference, c) {
				c.Tags = append(c.Tags, conference.Tags...)
				shouldAdd = false
			}
		}
		if shouldAdd {
			conference.Metadata = html_parser.GetHTMLMeta(conference.Url)
			filteredConferences = append(filteredConferences, conference)
			fmt.Println(conference)
		}
	}
	return filteredConferences
}
// the compare function returns true if "name" and "startDate" field is equal
func compare(a, b *Conference) bool {
	// log.Println("found duplicate", a.Name == b.Name && a.StartDate == b.StartDate)
	return a.Name == b.Name && a.StartDate == b.StartDate
}
