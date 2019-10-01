package main

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/yousseffarkhani/gophercises/09-quiet_hn/hn"
)

func main() {
	numStories := 30

	tmpl := template.Must(template.ParseFiles("index.gohtml"))

	http.HandleFunc("/", handler(numStories, tmpl))

	http.ListenAndServe(":8080", nil)

}

func handler(numStories int, tmpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, Data{
			Time:    time.Since(startTime),
			Stories: stories,
		})
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	numStories int
	cache      []item
	mutex      sync.Mutex
	duration   time.Duration
	expiration time.Time
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Now().Sub(sc.expiration) > 0 {
		stories, err := getTopStories(sc.numStories)
		if err != nil {
			return nil, err
		}
		sc.cache = stories
		sc.expiration = time.Now().Add(10 * time.Second)
		return sc.cache, nil
	}
	return sc.cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client

	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result)

	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{err: err}
			}
			resultCh <- result{item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item

	for _, res := range results {
		if isStoryLink(res.item) {
			if res.err != nil {
				continue
			}
			stories = append(stories, res.item)
		}
	}

	return stories
}

type item struct {
	hn.Item
	Host string
}

type Data struct {
	Stories []item
	Time    time.Duration
}

func parseHNItem(HNItem hn.Item) item {
	ret := item{Item: HNItem}
	url, err := url.Parse(ret.Url)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

func isStoryLink(parsedItem item) bool {
	return parsedItem.Type == "story" && parsedItem.Url != ""
}
