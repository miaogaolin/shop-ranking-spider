package main

import (
	"bytes"
	"log"
	url2 "net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Rakuten(content []byte) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	appendContent := regexp.MustCompile(`\/\*((\s|.)*?)\*\/`).FindAllStringSubmatch(string(content), -1)
	if len(appendContent) > 0 {
		for _, v := range appendContent {
			if !strings.Contains(v[1], "rnkRanking_after4box") {
				continue
			}
			doc.Find("#rnkRankingMain").AppendHtml(v[1])
		}
	}

	updateTime := regexp.MustCompile(`\d+年\d+月\d+日`).FindString(doc.Find("#rnkGenreRanking_updateDate").Text())
	var data []ProductList
	doc.Find("#rnkRankingMain>.rnkRanking_top3box,#rnkRankingMain>.rnkRanking_after4box").Each(func(i int, s *goquery.Selection) {
		img, _ := s.Find(".rnkRanking_imageBox img").Attr("src")
		title, _ := s.Find(".rnkRanking_imageBox img").Attr("alt")
		url, _ := s.Find(".rnkRanking_imageBox a").Attr("href")

		price := regexp.MustCompile(`\d+(,\d+)?`).FindString(s.Find(".rnkRanking_price").Text())

		priceBody, _ := strconv.ParseFloat(strings.Replace(price, ",", "", 1), 64)
		priceYen := "円"

		var productId string
		parseUrl, _ := url2.Parse(url)
		productId = endId(parseUrl.Path)

		shopUrl, _ := s.Find(".rnkRanking_shop a").Attr("href")
		storeId := endId(shopUrl)

		if priceBody == 0 {
			log.Printf("[Warn]产品价格为0, title=%s, url=%s", title, url)
			return
		}

		on := s.Find(".rnkRanking_starON").Length()
		half := s.Find(".rnkRanking_starHALF").Length()
		reviewRanking := float64(on) + float64(half)*0.5

		reviews := regexp.MustCompile(`\((\d+,?\d*)件\)`).FindStringSubmatch(s.Find(".rnkRanking_upperbox").Text())
		var reviewCount uint64
		if len(reviews) == 2 {
			count := strings.Replace(reviews[1], ",", "", 1)
			reviewCount, _ = strconv.ParseUint(count, 10, 64)
		}

		numHtml := s.Find(".rnkRanking_dispRank").Text()
		if numHtml == "" {
			numHtml, _ = s.Find(".rnkRanking_rankIcon img").Attr("alt")
		}
		if numHtml == "" {
			numHtml = s.Find(".rnkRanking_dispRank_overHundred").Text()
		}
		num, _ := strconv.ParseUint(strings.TrimRight(numHtml, "位"), 10, 64)
		data = append(data, ProductList{
			ID:            productId,
			Url:           url,
			Ranking:       num,
			Image:         img,
			Title:         title,
			PriceBody:     priceBody,
			PriceYen:      priceYen,
			StoreUrl:      shopUrl,
			StoreName:     storeId,
			StoreId:       storeId,
			ReviewRanking: reviewRanking,
			ReviewCount:   reviewCount,
		})
	})

	return ProductPage{
		UpdateTime: updateTime,
		List:       data,
	}, nil
}

func endId(path string) string {
	path = strings.TrimRight(path, "/")
	res := strings.Split(path, "/")
	if len(res) > 0 {
		return res[len(res)-1]
	}

	return ""
}
