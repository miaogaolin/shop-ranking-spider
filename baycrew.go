package main

import (
	"bytes"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func Baycrew(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var list []interface{}
	doc.Find("li.item").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".thumb a").Attr("href")
		num, _ := strconv.ParseUint(s.Find(".ranking").Text(), 10, 64)
		img, _ := s.Find(".thumb img").Attr("src")
		brand := s.Find(".brand").Text()
		cnName := s.Find(".data .name").Text()
		price := s.Find(".price").Text()

		list = append(list, map[string]interface{}{
			"url":           url,
			"ranking":       num,
			"image":         img,
			"brand":         brand,
			"brand_cn_name": cnName,
			"price":         price,
		})

	})
	return list, nil
}
