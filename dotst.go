package main

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func DotSt(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var list []interface{}
	doc.Find(".ranking-content__list li>a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		num, _ := strconv.ParseUint(s.Find(".icon-rank").Text(), 10, 64)
		img, _ := s.Find(".item-ph img").Attr("src")

		price := regexp.MustCompile(`¥\d+,?\d*`).FindString(s.Find(".item-price").Text())
		title := s.Find(".item-name").Text()
		brand := s.Find(".item-brand-name").Text()
		list = append(list, map[string]interface{}{
			"url":     url,
			"ranking": num,
			"image":   img,
			"price":   price,
			"title":   title,
			"brand":   brand,
			"id":      endId(url),
		})

	})
	return list, nil
}
