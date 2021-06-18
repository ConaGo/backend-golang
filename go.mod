module conago.de/web-scraper

go 1.16

replace conago.de/myutil => ../myutil

require (
	conago.de/myutil v0.0.0-00010101000000-000000000000
	github.com/PuerkitoBio/goquery v1.6.1 // indirect
	github.com/antchfx/xmlquery v1.3.6 // indirect
	github.com/antchfx/xpath v1.1.11 // indirect
	github.com/gocolly/colly/v2 v2.1.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	go.mongodb.org/mongo-driver v1.5.3
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	google.golang.org/appengine v1.6.7 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.10
)
