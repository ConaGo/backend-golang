package sqlite

//This module is a wrapper around the database connection
//it allows the connection to be reused in different goroutines
import (
	"log"

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
    dbC, errC := gorm.Open(sqlite.Open("../data/test.db"), &gorm.Config{})
    if errC != nil {
        log.Fatal("Failed to init test db:", errC)
    }
    ConfDB = DBT{DB: dbC}
    dbT, errT := gorm.Open(sqlite.Open("../data/quizToken.db"), &gorm.Config{})
    if errT != nil {
        log.Fatal("Failed to init quizToken db:", errT)
    }
    TokenDB = DBT{DB: dbT}
}