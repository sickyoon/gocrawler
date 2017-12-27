package main

import (
	"io"
	"log"

	"github.com/sickyoon/gocrawler/crawler"
	xmlpath "gopkg.in/xmlpath.v2"
)

// Request -
type Request struct {
	url    string
	method string
	page   int
}

// New -
func New(url string, page int) Request {
	return Request{
		url:    url,
		method: "GET",
		page:   page,
	}
}

// URL -
func (r Request) URL() string {
	return r.url
}

// Method -
func (r Request) Method() string {
	return r.method
}

// Process -
func (r Request) Process(reqCh chan crawler.Request, reader io.Reader) error {

	log.Printf("\n============================ page %d ============================\n", r.page)

	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		return err
	}
	it := xmlpath.MustCompile("//div[@class='quote']").Iter(root)
	for it.Next() {
		n := it.Node()
		text, _ := xmlpath.MustCompile("./span[@class='text']/text()").String(n)
		author, _ := xmlpath.MustCompile(".//small[@class='author']/text()").String(n)
		log.Printf("text: %s, author: %s", text, author)
		//tags, ok := xmlpath.MustCompile(".//div[@class='tags']/a[@class='tag']/text()").Iter(n)
	}
	nextURL, ok := xmlpath.MustCompile("//li[@class='next']/a/@href").String(root)
	// send request back
	if nextURL != "" && ok {
		reqCh <- New(nextURL, r.page+1)
	}

	return nil
}

func main() {

	// create scheduler
	s := crawler.New(1, 30)
	// start all workers
	err := s.StartWorkers()
	if err != nil {
		log.Fatalf("failed to create all workers")
	}

	// send request
	req := New("http://quotes.toscrape.com/", 1)
	s.Schedule(req)

	// TODO: should wait until all operations finish
	for {

	}
}
