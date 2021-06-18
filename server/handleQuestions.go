package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"database/sql/driver"

	"conago.de/web-scraper/db/mongo"
	"conago.de/web-scraper/db/sqlite"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)


func HandleQuestions(w http.ResponseWriter, r *http.Request) {
	//CORS - remove
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	path := r.URL.Path
	fmt.Println(path)
	params := r.URL.Query()
	amount := params.Get("amount")
	fmt.Println(amount)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
	mongo.QuizDB.DB.Find(ctx, bson.D{})
}
type Token struct {
	gorm.Model		`json:"-"`
	token string
	questionIDs []Question `gorm:"many2many:token_questions;"`
}
type Question struct {
	gorm.Model		`json:"-"`
	MongoI string
}
type mongoID primitive.ObjectID
func (id *mongoID) Scan(src interface{}) error {
	if t, ok := src.(primitive.ObjectID); ok {
		*id = mongoID(t)
	}
	return nil
}
func (id mongoID) Value() (driver.Value, error) {
	return fmt.Sprint(id), nil
}
func Test() {
	//sqlite.TokenDB.DB.AutoMigrate((&Token{}))

	sqlite.TokenDB.DB.AutoMigrate((&Question{}))
	//var q = []question{{MongoI: "jklsad"}, {MongoI: "hjk"}}
	//sqlite.TokenDB.DB.Create(&q)
	var ctx = context.TODO()

	//FindALL
/* 	cursor, err := mongo.QuizDB.DB.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return 
	}
	var questions []bson.M 
	if err = cursor.All(ctx, &questions); err != nil {
		log.Fatal(err)
	} 
		fmt.Println((questions))
	*/
	var question bson.M
	if err := mongo.QuizDB.DB.FindOne(ctx, bson.M{}).Decode(&question); err != nil {
		log.Fatal(err)
	}
	fmt.Println(question["_id"])
	fmt.Println("_________")
	//idq := question["_id"].(primitive.ObjectID)
	xType := reflect.TypeOf(fmt.Sprint(question["_id"]))
xValue := reflect.ValueOf(fmt.Sprint(question["_id"]))

fmt.Println(xType, xValue) 

	fmt.Println(fmt.Sprint(question["_id"]))
	q := Question{MongoI: fmt.Sprint(question["_id"])}
	fmt.Println(q.MongoI)
	sqlite.TokenDB.DB.Create(&q)

	questions := []Question{}
	sqlite.TokenDB.DB.Find(&questions)
	fmt.Println(questions)




/* 	db.AutoMigrate(&Conference{})
	db.AutoMigrate((&Tag{}))
	datas := readConferenceData()
	db.Create(&datas)
	var conferences []Conference
	for _, elem := range conferences {
		fmt.Println(elem.Name, elem.City)
	} */
}



func HandleTokens(w http.ResponseWriter, r *http.Request){
	id := uuid.New()
	fmt.Fprint(w, id.String())
}

