package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL      string   `xml:"guid"`
	Title    string   `xml:"title"`
	Category []string `xml:"category"`
}

func FetchContent(url string) (*RSS, error) {
	fmt.Printf("fetching: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET for URL %s: %w", url, err)
	}

	defer resp.Body.Close()

	rss := new(RSS)
	if err = xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("unmarshal body: %w", err)
	}

	return rss, nil
}

//type List struct {
//	Items []Tags
//}
//
//type Tags struct {
//	URL      string
//	Title    string
//	Category []string
//}
//
//// обновление кеша по одному листу, например, по go через представление List
//func FetchContent(url string) (*List, error) {
//	fmt.Printf("fetching: %s\n", url)
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, fmt.Errorf("HTTP GET for URL %s: %w", url, err)
//	}
//	defer resp.Body.Close()
//
//	list := new(List)
//	if err = xml.NewDecoder(resp.Body).Decode(&list); err != nil {
//		return nil, fmt.Errorf("unmarshal body: %w", err)
//	}
//
//	return list, nil
//}
