package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
)

type SearchResult struct {
	Description       string `json:"description"`
	DisplayPrice      string `json:"displayPrice"`
	ProductFamilyName string `json:"productFamilyName"`
	ProductId         string `json:"productId"`
	PublisherName     string `json:"publisherName"`
	Title             string `json:"title"`
}

func searchStore(q string) ([]SearchResult, error) {
	// https://apps.microsoft.com/api/products/search?gl=US&hl=en-us&query=calculator&mediaType=all&age=all&price=all&category=all&subscription=all&cursor=

	resp, err := http.Get("https://apps.microsoft.com/api/products/search?gl=US&hl=en-us&query=" + url.QueryEscape(q) + "&mediaType=all&age=all&price=all&category=all&subscription=all&cursor=")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var results []SearchResult
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		result := SearchResult{}
		result.Description, _ = jsonparser.GetString(value, "description")
		result.DisplayPrice, _ = jsonparser.GetString(value, "displayPrice")
		result.ProductFamilyName, _ = jsonparser.GetString(value, "productFamilyName")
		result.ProductId, _ = jsonparser.GetString(value, "productId")
		result.PublisherName, _ = jsonparser.GetString(value, "publisherName")
		result.Title, _ = jsonparser.GetString(value, "title")

		results = append(results, result)

	}, "productsList")

	if len(results) == 0 {
		return nil, errors.New("no results found")
	}
	if len(results) > 5 {
		results = results[:5]
	}

	return results, nil
}

type InstallationURI struct {
	Name string
	URI  string
}

func getFiles(prodId string) ([]InstallationURI, error) {
	resp, err := http.Post(
		"https://store.rg-adguard.net/api/GetFiles",
		"application/x-www-form-urlencoded",
		strings.NewReader("type=ProductId&url="+prodId+"&ring=RP&lang=en-US"),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []InstallationURI
	doc.Find("td > a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		if (strings.HasSuffix(name, ".appx") || strings.HasSuffix(name, ".appxbundle")) &&
			(strings.Contains(name, "x64") || strings.Contains(name, "neutral")) {
				
			files = append(files, InstallationURI{
				Name: s.Text(),
				URI:  s.AttrOr("href", ""),
			})
		}
	})

	return files, nil
}
