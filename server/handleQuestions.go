package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"conago.de/myutil"
	mongodb "conago.de/web-scraper/db/mongo"
	sqlitedb "conago.de/web-scraper/db/sqlite"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)
func init() {
	sqlitedb.TokenDB.DB.AutoMigrate((&Token{}))
	sqlitedb.TokenDB.DB.AutoMigrate((&Question{}))
}

func HandleQuestions(w http.ResponseWriter, r *http.Request) {
	//CORS - remove
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//parse parametes amount & token & difficulty
	params := r.URL.Query()
	amount, err := strconv.Atoi(params.Get("amount"))
	if err != nil || amount < 1 || amount > 50 {
		amount = 10
	}
	tokenParam := params.Get("token")
	difficulty := params.Get("difficulty")
	if !myutil.StringInSlice(difficulty, []string{"easy", "medium", "hard"}) {
		difficulty = ""
	}
	//if a token is sent, prepare the query to filter questions that were already sent to this client
	amountStage :=  bson.D{{"$sample",bson.D{{"size",amount}}}}
	pipeline := mongo.Pipeline{amountStage}
	if tokenParam != ""{
		tok := Token{}
		result := sqlitedb.TokenDB.DB.Preload("QuestionIDs").Where("token = ?", tokenParam).First(&tok)
		//fmt.Println(tok.QuestionIDs)
		oids := make([]primitive.ObjectID, len(tok.QuestionIDs))
		if result.Error == nil {
			
			for i := range tok.QuestionIDs {
			  oids[i], _ = primitive.ObjectIDFromHex(tok.QuestionIDs[i].MongoID)
			}
			tokenStage := bson.D{{"$match",bson.D{{"_id", bson.D{{"$nin",oids }}}}}}
			pipeline = mongo.Pipeline{tokenStage, amountStage}
			fmt.Println("added")
		} 
	}	

	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := mongodb.QuizDB.DB.Aggregate(ctx, pipeline, opts)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	if tokenParam != "" {
		var ids []string
		for _, result := range results {
			fmt.Printf("name %v\n", result["_id"])
			ids = append(ids, result["_id"].(primitive.ObjectID).Hex())	
		}
		var token = Token{}
		var questions []Question
		sqlitedb.TokenDB.DB.First(&token, "token = ?", tokenParam)
		sqlitedb.TokenDB.DB.Find(&questions, "mongo_id IN (?)", ids)
		sqlitedb.TokenDB.DB.Model(&token).Association("QuestionIDs").Append(&questions)		
	}
	response, err := json.Marshal(results)
	w.Write(response)
}
type Token struct {
	gorm.Model		`json:"-"`
	Token string
	QuestionIDs []Question `gorm:"many2many:token_questions;"`
}
type Question struct {
	gorm.Model		`json:"-"`
	MongoID string
}
func Test() {
	sqlitedb.TokenDB.DB.AutoMigrate((&Token{}))
	sqlitedb.TokenDB.DB.AutoMigrate((&Question{}))


	var ctx = context.TODO()

	cursor, err := mongodb.QuizDB.DB.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var questions []Question
	for cursor.Next(ctx) {
		var question bson.M
		if err = cursor.Decode(&question); err != nil {
			log.Fatal(err)
		}
		mongoID := question["_id"].(primitive.ObjectID).Hex()
		q := Question{MongoID: mongoID}
		questions = append(questions, q)
		
	}
	sqlitedb.TokenDB.DB.Create(&questions)
}
func HandleTokens(w http.ResponseWriter, r *http.Request){
	id := uuid.New()
	sqlitedb.TokenDB.DB.Create(&Token{Token:id.String()})
	fmt.Fprint(w, id.String())
}

