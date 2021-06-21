package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	MongoID string
}
func Test() {
	sqlite.TokenDB.DB.AutoMigrate((&Token{}))
	sqlite.TokenDB.DB.AutoMigrate((&Question{}))

	
	var ctx = context.TODO()

	var question bson.M
	if err := mongo.QuizDB.DB.FindOne(ctx, bson.M{}).Decode(&question); err != nil {
		log.Fatal(err)
	}
	mongoID := question["_id"].(primitive.ObjectID).Hex()
	q := Question{MongoID: mongoID}
	sqlite.TokenDB.DB.Create(&q)

	questions := []Question{}
	sqlite.TokenDB.DB.Find(&questions)
	fmt.Println(questions)
}



func HandleTokens(w http.ResponseWriter, r *http.Request){
	id := uuid.New()
	fmt.Fprint(w, id.String())
}

