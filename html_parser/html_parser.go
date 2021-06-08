package html_parser

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type HTMLMeta struct {
	Title         string
	Description   string
	OGTitle       string
	OGDescription string
	OGImage       string
	OGAuthor      string
	OGPublisher   string
	OGSiteName    string
}

func GetHTMLMeta(url string) HTMLMeta {
	response, err := http.Get(url)
	var hm HTMLMeta
	if err != nil {
		log.Printf("Error getting HTML from url %s ::ERROR:: %s", url, err)
	} else {
		defer response.Body.Close()
		hm = extractMetaData(response.Body)
	}
	return hm
}
func extractMetaData(resp io.Reader) (hm HTMLMeta) {
	z := html.NewTokenizer(resp)

	titleFound := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				//fmt.Println(t.Attr)
				if desc, ok := extractMetaProperty(t, "description"); ok {
					hm.Description = desc
				}
				if ogTitle, ok := extractMetaProperty(t, "og:title"); ok {
					hm.OGTitle = ogTitle
				}
				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					hm.OGDescription = ogDesc
				}

				ogImage, ok := extractMetaProperty(t, "og:image")
				if ok {
					hm.OGImage = ogImage
				}

				ogAuthor, ok := extractMetaProperty(t, "og:author")
				if ok {
					hm.OGAuthor = ogAuthor
				}

				ogPublisher, ok := extractMetaProperty(t, "og:publisher")
				if ok {
					hm.OGPublisher = ogPublisher
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					hm.OGSiteName = ogSiteName
				}
			}
		case html.TextToken:
			if titleFound {
				t := z.Token()
				hm.Title = t.Data
				titleFound = false
			}
		}

	}
}
func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}
		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}
