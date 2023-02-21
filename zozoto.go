package main

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Zozoto(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var list []interface{}
	doc.Find(".catalog-item-container").Each(func(i int, s *goquery.Selection) {
		num := s.Find(".catalog-hero>div").Eq(1).Text()
		img, _ := s.Find("img").Attr("src")
		title := s.Find(".catalog-property").Text()
		url, _ := s.Find(".catalog-link").Attr("href")

		price := s.Find(".catalog-price-number").Text()
		brand := strings.TrimSpace(s.Find(".catalog-header-h").Text())
		list = append(list, map[string]interface{}{
			"url":     url,
			"ranking": num,
			"image":   img,
			"title":   title,
			"price":   price,
			"brand":   brand,
		})

	})
	return list, nil
}
