package main

import (
	//"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
)

func GetDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url) 
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

type Piece struct {
	url string
	title string
	show *Show
}

type Show struct {
	url string
	doc *goquery.Document
}

func GetShow(url string) *Show {
	doc := GetDocument(url)
	show := Show{
		url: url,
		doc: doc,
	}
	return &show
}

func ParsePieceSelection(s *goquery.Selection) *Piece {
	showHead := s.Find(".views-group-header")
	if (showHead.Length() != 0) {
		url := showHead.Find("")
		show := &Show{
			url: url,
		}
	}
	piece := &Piece{
		title: s.Find(".node-title").Text(),
	}
	return piece
}

func DocPieces(doc *goquery.Document) []*Piece {
	groups := doc.Find(".view-stories .views-group")
	pieces := make([]*Piece, groups.Length())
	groups.Each(func(i int, s *goquery.Selection) {
		pieces[i] = ParsePieceSelection(s)
  	})
	return pieces
}

func Scrape() []*Piece {
	url := "http://www.thestory.org/stories"
	pieces := DocPieces(GetDocument(url))
	return pieces
}

func main() {
	Scrape()
}

func doLog() {
	// fmt.Printf("Review %d: %s - %s\n", i, band, title)
	// fmt.Printf("Episode: %s\n", title)
}
