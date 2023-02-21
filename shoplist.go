package main

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Shoplist(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	var list []interface{}
	doc.Find("ol.p-list_product>li").Each(func(i int, s *goquery.Selection) {
		num, _ := s.Find(".p-list_product_link").Attr("data-index")
		img, _ := s.Find(".p-list_product_imgwrap img").Attr("src")
		title := s.Find(".op_title").Text()
		score := strings.Split(s.Find(".p-list_product_review_score").Text(), "（")
		var reviewRanking, reviewCount string
		if len(score) == 2 {
			reviewRanking = score[0]
			reviewCount = strings.TrimRight(score[1], "）")
		}
		url, _ := s.Find(".p-list_product_link").Attr("href")

		price := regexp.MustCompile(`¥\d+,?\d*`).FindString(s.Find(".p-list_product_pricewrap").Text())
		brand := strings.TrimSpace(s.Find(".p-list_product_txtwrap .u-txt_rdstr:first-child").Text())
		list = append(list, map[string]interface{}{
			"url":            url,
			"ranking":        num,
			"image":          img,
			"title":          title,
			"review_ranking": reviewRanking,
			"review_count":   reviewCount,
			"price":          price,
			"brand":          brand,
		})

	})
	return list, nil
}
