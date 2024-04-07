package main

import (
	"flag"
	"fmt"
	"io"
	"link"
	linklib "link"
	"log"
	"net/http"
	urlpkg "net/url"
)

/**
findAllLinks that takes a url and finds all links accessible from that url to a max depth
 - list of links generated for landing url (flag set by user on command line)
 - iterate over list of links and do findAllLinks on each?
	- pass depth to this findAllLinks which returns if depth > depth

can use a map as a set in go. set of all urls makes sure there are no duplicates and acts
as a check to say 'have i visited this site already'

can i utilize the link package from previous exercises? this package builds a slice of all
anchors and their corresponding hrefs. problem is, it relies on the presence of a local html
file and does file parsing vs network fetching. might be worth writing an adapted version of that
algorithm just for this project.
**/

var maxDepth int
var url string
var baseURL string

func main() {
	u, d := parseFlags()
	url = u
	baseURL = getURLRoot(url)
	maxDepth = d
	allLinks := findAllLinks(make([]linklib.Link, 0), make(map[linklib.Link]bool), 0)
	fmt.Println(allLinks)
	// sitemap = writeXML(allLinks)
	// fmt.Println(sitemap)
}

func parseFlags() (string, int) {
	var url = flag.String("url", "https://github.com", "URL for the site to build sitemap from")
	var depth = flag.Int("depth", 4, "maximum recursive depth")
	flag.Parse()
	return *url, *depth
}

func getHTML(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Bad response fetching url %s\n", url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Couldn't read response")
	}
	return body
}

func getURLRoot(urlStr string) string {
	u, err := urlpkg.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	fmt.Println(u.Host)
	return u.String()
}

func findAllLinks(links []linklib.Link, seen map[linklib.Link]bool, currdepth int) []link.Link {
	if currdepth > maxDepth {
		return links
	}
	if currdepth == 0 {
		htmlBuff := getHTML(url)
		links = link.Parse(htmlBuff)
	}
	for _, link := range links {
		if haveSeen := seen[link]; haveSeen {
			continue
		}
		seen[link] = true
		if getURLRoot(link.Href) == baseURL {
			html := getHTML(link.Href)
			links = append(links, findAllLinks(linklib.Parse(html), seen, currdepth+1)...)
		}
	}

	return links
}

// func writeXML()
// <?xml version="1.0" encoding="UTF-8"?>
// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//   <url>
//     <loc>http://www.example.com/</loc>
//   </url>
//   <url>
//     <loc>http://www.example.com/dogs</loc>
//   </url>
// </urlset>
