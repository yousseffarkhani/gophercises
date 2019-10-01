package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBase = "https://hacker-news.firebaseio.com/v0"

type Client struct {
	apiBase string
}

func (client *Client) defaultify() {
	if client.apiBase == "" {
		client.apiBase = apiBase
	}
}

func (client *Client) TopItems() ([]int, error) {
	client.defaultify()
	var ids []int

	resp, err := http.Get(fmt.Sprintf("%s/topstories.json", client.apiBase))
	if err != nil {
		return ids, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

func (client *Client) GetItem(id int) (Item, error) {
	client.defaultify()
	var item Item

	resp, err := http.Get(fmt.Sprintf("%s/item/%d.json", client.apiBase, id))
	if err != nil {
		return item, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	Text string `json:"text"`
	Url  string `json:"url"`
}
