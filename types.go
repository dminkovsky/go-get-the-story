package main

import (
	"fmt"
	//"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mxk/go-sqlite/sqlite3"
)

type ListPage struct {
	doc      *goquery.Document
	path     string
	nextPath string
	segments []*Segment
}

func ParseListPageHead(head *goquery.Selection) *Show {
	path, _ := head.Find("h4 a").Attr("href")
	return &Show{path: path}
}

func ParseListPageSeg(seg *goquery.Selection) *Segment {
	a := seg.Find(".node-title a")
	path, _ := a.Attr("href")
	title := a.Text()
	return &Segment{path: path, title: title}
}

func (p *ListPage) Scrape() {
	p.doc = GetDocument(GetUrl(p.path))
	segs := p.doc.Find(".view-stories .views-row")
	p.segments = make([]*Segment, segs.Length())
	segs.Each(func(i int, seg *goquery.Selection) {
		segment := ParseListPageSeg(seg)
		head := seg.Closest(".views-group-inner").Find(".views-group-header")
		if head.Length() != 0 {
			segment.show = ParseListPageHead(head)
		}
		p.segments[i] = segment
	})
}

func (p *ListPage) Next() string {
	next := p.doc.Find("li.pager-next a")
	if next.Length() == 0 {
		return ""
	}
	path, _ := next.Attr("href")
	return GetUrl(path)
}

type Segment struct {
	path  string
	title string
	show  *Show
}

func (s *Segment) Scrape() {}

func (s *Segment) Save(conn *sqlite3.Conn) {}

type Show struct {
	path string
	doc  *goquery.Document
}

func (s *Show) Exists(conn *sqlite3.Conn) bool {
	var count int64
	stmt, err := conn.Query(`SELECT count(*) FROM shows WHERE path="$p"`, sqlite3.NamedArgs{"$p": s.path})
	logFatal(err)
	stmt.Scan(&count)
	stmt.Close()
	return count == 1
}

func (s *Show) Scrape() {
	doc := GetDocument(GetUrl(s.path))
	fmt.Println(s.path, doc)
}

func (s *Show) Save(conn *sqlite3.Conn) {}
