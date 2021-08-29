package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// struct to store url and text
type Link struct {
	url  string
	text string
}

func main() {
	args := os.Args
	log.Printf("args received %+v \n", args)
	extractor("https://bhanureddy.dev")
}

/*
extractor moves from one link to another in a given page
*/
func extractor(url string) {
	resp, err := http.Get(url) // send request and get response
	if err != nil {
		log.Printf("error occured while getting response http.get(%s) and error is %s", url, err.Error())
	}
	defer resp.Body.Close()

	links, err := collectLinksFromHtml(resp.Body)
	if err != nil {
		log.Fatal("failed collecting links")
	}

	for _, link := range links {
		nonDuplicateLinks := make(map[string]int, 1)
		if _, ok := nonDuplicateLinks[link.url]; !ok {
			nonDuplicateLinks[link.url] = 1
			log.Printf("URL: %s TEXT: %s \n", link.url, link.text)
		}
	}
}

/*
collectLinksFromHtml reads htmls response from request parses html page
constructs html nodes  and navigates through all html nodes and extracts
a tag elements and then actual urls
*/
func collectLinksFromHtml(htmlResp io.Reader) ([]Link, error) {
	doc, err := html.Parse(htmlResp)
	if err != nil {
		log.Printf("error while parsing htmlResp io.Reader in html.Parse(), error: %s", err.Error())
		return nil, err
	}
	//log.Printf("htmlNOde %+v \n", doc)
	nodes := getHtmlNodes(doc)
	//log.Printf("nodes %+v \n", nodes)

	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

/*
getHtmlNodes collects html node elements of tag type a
uses recursion to traverse through all the a tags available in the page
and returns them
*/
func getHtmlNodes(doc *html.Node) []*html.Node {
	if doc.Type == html.ElementNode && doc.Data == "a" {
		log.Printf("tag a %+v \n", doc)
		return []*html.Node{doc}
	}

	var nod []*html.Node
	for b := doc.FirstChild; b != nil; b = b.NextSibling {
		nod = append(nod, getHtmlNodes(b)...)
	}
	return nod
}

/*
buildLink extracts actual url and associate text
*/
func buildLink(node *html.Node) Link {
	var link Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.url = attr.Val
		}
	}

	link.text = getText(node)
	return link
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return strings.Join(strings.Fields(text), " ")
}
