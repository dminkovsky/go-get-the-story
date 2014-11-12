package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/mxk/go-sqlite/sqlite3"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var baseUrl = "http://www.thestory.org"
var conn = GetSqliteConn("go-get-the-story.db")

func GetUrl(path string) string {
	return baseUrl + path
}

func GetFile(url string) []byte {
	res, err := http.Get(url)
	logFatal(err)
	defer res.Body.Close()
	file, err := ioutil.ReadAll(res.Body)
	logFatal(err)
	return file
}

func GetDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	logFatal(err)
	return doc
}

func GetSqliteConn(db string) *sqlite3.Conn {
	conn, err := sqlite3.Open(db)
	logFatal(err)
	return conn
}

func InitDb(conn *sqlite3.Conn) {
	conn.Exec(`CREATE TABLE IF NOT EXISTS shows(
		id         INTEGER PRIMARY KEY,
		path       NOT NULL,
		title      NOT NULL,
		date       NOT NULL
	);`)

	conn.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS showPath ON shows(path);`)

	conn.Exec(`CREATE TABLE IF NOT EXISTS segments(
		id         INTEGER PRIMARY KEY,
		path       NOT NULL,
		title      NOT NULL,
		date       NOT NULL,
		show_id  
	);`)

	conn.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS segmentPath ON segments(path);`)
}

func Crawl(path string) {
	page := ListPage{path: path}
	page.Scrape()

	for _, segment := range page.segments {
		if segment.show != nil {
			show := segment.show
			if !show.Exists(conn) {
				show.Scrape()
				show.Save(conn)
			}
		}
		segment.Scrape()
		segment.Save(conn)
	}

	next := page.Next()
	if next != "" {
		Crawl(next)
	}
}

func main() {
	fmt.Println("go-get-the-story\n")

	InitDb(conn)

	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "/stories"
	}

	Crawl(path)
}

func doLog() {
	// fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// fmt.Printf("Episode: %s\n", title)
}
