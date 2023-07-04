package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string){
	for _,a := range t.Attr{
		if a.Key == "href"{
			href = a.Val
			ok = true
		}
	}
	return
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		chFinished <- true
	}()
	if err != nil {
		fmt.Println("ERROR:failed to crawl: ", url)
		return
	}
	b := resp.Body

	defer b.Close()

z:=html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)

			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}

	}
}

func main() {

	foundUrl := make(map[string]bool)
	seedUrl := os.Args[1:]

	chUrl := make(chan string)
	chFinished := make(chan bool)

	for _, url := range seedUrl {
		go crawl(url, chUrl, chFinished)
	}
	for c := 0; c < len(seedUrl); {
		select {
		case url := <-chUrl:
			foundUrl[url] = true
		case <-chFinished:
			c++
		}
	}
	fmt.Println("\nFound", len(foundUrl), "unique url: \n")

	for url, _ := range foundUrl {
		fmt.Println("-" + url)
	}
	close(chUrl)
}
