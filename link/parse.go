package link

import (
	"bytes"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string `xml:"loc"`
	Text string `xml:"text"`
}

func Parse(b []byte) []Link {
	head := parseFile(b)
	linkNodes := GetAllLinkNodes(head)
	var result []Link
	for _, node := range linkNodes {
		result = append(result, Link{Href: node.Attr[0].Val, Text: buildText(node)})
	}
	return result
}

func parseFile(file []byte) *html.Node {
	r := bytes.NewReader(file)
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal("Couldn't parse file", err)
	}
	return doc
}

func buildText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += buildText(c)
	}
	return strings.Join(strings.Fields(text), " ")
}

func GetAllLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, GetAllLinkNodes(c)...)
	}
	return nodes
}
