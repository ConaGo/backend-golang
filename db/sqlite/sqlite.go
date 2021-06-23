package sqlite

//This module is a wrapper around the database connection
//it allows the connection to be reused in different goroutines
import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
type DBT struct {
    DB *gorm.DB
}

var ConfDB DBT
var TokenDB DBT

//opens two connections to the sqlite databases and exposes the connection
func init() {
    godotenv.Load(".env")
    dbC, errC := gorm.Open(sqlite.Open(os.Getenv("CONFERENCEDB_PATH")), &gorm.Config{})
    if errC != nil {
        log.Fatal("Failed to init test db:", errC)
    }
    ConfDB = DBT{DB: dbC}
    dbT, errT := gorm.Open(sqlite.Open(os.Getenv("TOKENDB_PATH")), &gorm.Config{})
    if errT != nil {
        log.Fatal("Failed to init quizToken db:", errT)
    }
    TokenDB = DBT{DB: dbT}
}