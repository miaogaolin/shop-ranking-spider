package main

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	url2 "net/url"

	"github.com/PuerkitoBio/goquery"
)

func Yahoo(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	updateTime := regexp.MustCompile(`\d+/\d+/\d+`).FindString(doc.Find(".NBOPF25lQ6VT").Text())
	var list []ProductList
	doc.Find(".list .line").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find(".name a").Attr("href")
		num, _ := strconv.ParseUint(s.Find(".rank-text").Text(), 10, 64)
		img, _ := s.Find(".column-middle-left .image").Attr("src")
		title := s.Find(".name-text").Text()
		price := s.Find(".price-number").Text()
		priceBody, _ := strconv.ParseFloat(strings.Replace(price, ",", "", 1), 64)
		shopUrl, _ := s.Find(".store-name .store-link").Attr("href")
		storeName := s.Find(".store-text").Text()
		reviewRanking, _ := strconv.ParseFloat(s.Find(".review-average").Text(), 64)
		countHtml := s.Find(".review-count").Text()
		countStr := strings.Trim(regexp.MustCompile(`\d+,?\d*件`).FindString(countHtml), "件")
		countStr = strings.Replace(countStr, ",", "", -1)
		reviewCount, _ := strconv.ParseUint(countStr, 10, 64)
		list = append(list, ProductList{
			ID:            ProductID(url),
			Url:           url,
			Ranking:       num,
			Image:         img,
			Title:         title,
			PriceBody:     priceBody,
			StoreUrl:      shopUrl,
			StoreName:     storeName,
			StoreId:       StoreID(shopUrl),
			ReviewRanking: reviewRanking,
			ReviewCount:   reviewCount,
		})
	})
	return ProductPage{
		UpdateTime: updateTime,
		List:       list,
	}, nil
}

func ProductID(url string) string {
	u, _ := url2.Parse(url)
	paths := strings.Split(u.Path, "/")
	if len(paths) > 0 && strings.Contains(paths[len(paths)-1], "html") {
		p := strings.Split(paths[len(paths)-1], ".")
		if len(p) == 2 {
			return p[0]
		}
	}
	return ""
}

func StoreID(url string) string {
	u, _ := url2.Parse(url)
	return endId(u.Path)
}
