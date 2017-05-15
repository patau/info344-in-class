package main

import (
	"fmt"
	"os"
	"time"
)

const usage = `
usage:
	crawler <starting-url>
`

//linkq: links waiting to be fetched
//resultsq: all the hyperlinks that we fetched
func worker(linkq chan string, resultsq chan []string) {
	for link := range linkq { //for all the links
		plinks, err := getPageLinks(link)
		if err != nil {
			fmt.Printf("ERROR fetching %s: %v\n", link, err)
			//exit the current loop itr and go to the next one
			//go back to the main for loop, grab the next item in q
			continue
		}
		fmt.Printf("%s (%d links)\n", link, len(plinks.Links))
		time.Sleep(time.Millisecond * 500)
		if len(plinks.Links) > 0 {
			//Anonymous go func
			go func(links []string) {
				resultsq <- links
			}(plinks.Links)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	nWorkers := 50
	linkq := make(chan string, 1000)
	resultsq := make(chan []string, 1000)
	for i := 0; i < nWorkers; i++ {
		go worker(linkq, resultsq)
	}

	linkq <- os.Args[1] //the starting url

	seen := map[string]bool{}
	for links := range resultsq { //Read out of the results q; process links inside
		for _, link := range links {
			if !seen[link] {
				seen[link] = true
				linkq <- link
			}
		}
	}

}
