package quiz_scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Category struct {
	Id 			int				`gorm:"primaryKey"`
	Name 		string			`gorm:"index"`
	CreatedAt 	time.Time
  	UpdatedAt 	time.Time
  	DeletedAt 	gorm.DeletedAt 
}
type Question struct {
	ID     			primitive.ObjectID 	`bson:"_id,omitempty"`
	Category		string
	Type 			string
	Difficulty		string
	Question		string
	CorrectAnswer	string				`json:"correct_answer" bson:"correct_answer,omitempty"`
	IncorrectAnswer	[]string			`json:"incorrect_answers" bson:"incorrect_answers,omitempty"`
}
type Response struct {
	ResponseCode	int 				`bson:"response_code"`
	Results 		[]Question
}
func GetQuizzes() {
	godotenv.Load("../.env")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database("questions")
	collection := database.Collection("question")
	amount := 50
	for amount > 0 {
		resp, httpErr := http.Get(getUrl(amount))
		if httpErr != nil {
			// handle error
			fmt.Println(httpErr)
			return
		}
		defer resp.Body.Close()
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			fmt.Println(readErr)
			return
		}
		var questions []Question
		var response  Response
		parseErr := json.Unmarshal(body, &response)
		if parseErr != nil {
			fmt.Println(parseErr)
			return
		}
		questions = response.Results

		fmt.Println(len(questions))
		fmt.Println(response.ResponseCode)
		if response.ResponseCode > 0 {
			amount--
		}

		for _, e := range questions {
			res, insertErr := collection.InsertOne(ctx, e)
			if insertErr != nil {
					log.Fatal(insertErr)
			}
			fmt.Println(res);
		}
		

	}

}
func getUrl(amount int) string{
	token := "5a9fbc6952f21fb76813c0f7d2428d7d930c64bd5edf208c6657a3ee1f3a89d6"

	return fmt.Sprintf("https://opentdb.com/api.php?amount=%v&token=%s",amount, token)
}
