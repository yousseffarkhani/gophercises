/* Pour exporter le résultat du programme dans un fichier, utiliser go run main.go > map.xml */
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/yousseffarkhani/gophercises/04-linkParser/linkParser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com/demos/cyoa/", "url to scan")
	maxDepth := flag.Int("depth", 3, "Max depth to search")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)

	toXml := urlset{
		Xmlns: xmlns,
		Urls:  make([]loc, len(pages)),
	}
	for i, page := range pages {
		toXml.Urls[i] = loc{page}
	}

	fmt.Printf(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{}) // struct utilise moins d'espace en mémoire par rapport à un bool
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{} // Marks url as seen
			for _, link := range get(url) {
				if _, ok := seen[url]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}
	ret := make([]string, 0, len(seen)) // Optimisation du slice (permet d'allouer de la mémoire en avance de phase)
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func retrieveLinks(body io.Reader, base string) []string {
	links, err := linkParser.Parse(body)
	checkError(err)
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr) // Il faut absolument fermer la response.body
	checkError(err)
	html := resp.Body
	defer html.Close()
	// io.Copy(os.Stdout, html) // Permet d'afficher un io.Reader dans stdout

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{ // Intéressant pour récupérer l'URL de base
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(retrieveLinks(html, base), withPrefix(base))
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeDuplicates(urls []string) []string {
	var newList []string
	for _, url := range urls {
		if count(url, newList) == 0 {
			newList = append(newList, url)
		}
	}
	return newList
}

func count(s string, slice []string) int {
	var count int
	for _, word := range slice {
		if word == s {
			count++
		}
	}
	return count
}
