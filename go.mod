module conago.de/web-scraper

go 1.16

replace conago.de/myutil => ../myutil

require (
	conago.de/myutil v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	go.mongodb.org/mongo-driver v1.5.3
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.10
)
