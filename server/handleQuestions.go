package server

import (
	"context"
	"net/http"
	"time"

	"conago.de/web-scraper/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)


func HandleQuestions(w http.ResponseWriter, r *http.Request) {
	//CORS - remove
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
	mongo.QuizDB.DB.Find(ctx, bson.D{})
}
