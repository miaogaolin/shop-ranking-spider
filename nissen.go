package main

import (
	"bytes"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func Nissen(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var list []interface{}
	doc.Find(".rank_list>ul>li").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".m-item_list_txt a").Attr("href")
		title := s.Find(".m-item_list_txt span").Text()
		priceOutTax := s.Find(".item-price-out-tax").Text()
		priceInTax := s.Find(".item-price-in-tax").Text()
		reviewRanking := s.Find(".m-item_list_reviewscore").Text()
		num := s.Find(".m-item_list_rank").Text()
		img, _ := s.Find(".m-item_list_photo_main img").Attr("src")
		reviewCount := regexp.MustCompile(`\d+件`).FindString(s.Find(".m-item_list_review").Text())
		size := s.Find(".item-size").Text()
		list = append(list, map[string]interface{}{
			"url":            url,
			"ranking":        num,
			"image":          img,
			"price_out_tax":  priceOutTax,
			"price_in_tax":   priceInTax,
			"title":          title,
			"review_ranking": reviewRanking,
			"review_count":   reviewCount,
			"size":           size,
		})

	})
	return list, nil
}
