package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	linklib "link"
	"log"
	"net/http"
	urlpkg "net/url"
	"os"
)

var maxDepth int
var url string
var baseURL string

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	u, d := parseFlags()
	url = u
	baseURL = getURLRoot(url)
	maxDepth = d
	allLinks := buildMap(make([]linklib.Link, 0), make(map[linklib.Link]bool), 0)
	xml := toXml(allLinks)
	writeXML(xml)
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
		return ""
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	u.Scheme = ""
	return u.String()
}

func buildMap(links []linklib.Link, seen map[linklib.Link]bool, currdepth int) []linklib.Link {
	if currdepth > maxDepth {
		var empty []linklib.Link
		return empty
	}
	if currdepth == 0 {
		htmlBuff := getHTML(url)
		links = linklib.Parse(htmlBuff)
	}
	validLinks := make([]linklib.Link, 0)
	for _, link := range links {
		if haveSeen := seen[link]; haveSeen {
			continue
		}
		seen[link] = true
		if getURLRoot(link.Href) == baseURL {
			validLinks = append(validLinks, link)
			html := getHTML(link.Href)
			validLinks = append(validLinks, buildMap(linklib.Parse(html), seen, currdepth+1)...)
		}
	}
	return validLinks
}

func toXml(links []linklib.Link) urlset {
	set := urlset{Xmlns: xmlns}
	for _, link := range links {
		set.Urls = append(set.Urls, loc{Value: link.Href})
	}
	return set
}

func writeXML(urls urlset) {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(urls); err != nil {
		panic(err)
	}
	fmt.Println()
}
