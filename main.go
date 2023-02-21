package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("需要指定爬虫数据文件")
	}
	filename := os.Args[1]
	ext := filepath.Ext(filename)

	var (
		content []byte
		err     error
	)

	if ext == ".html" {
		content, err = os.ReadFile(filename)
	} else if ext == ".mhtml" {
		if strings.Contains(filename, "ZOZOTO") {
			IsJapanese = true
		}
		content, err = MthmlToHtml(filename)
	} else {
		log.Fatal("仅支持 .html、.mhtml 文件格式")
	}
	if err != nil {
		log.Fatal(err)
	}
	var res interface{}
	if strings.Contains(filename, "楽天") || strings.Contains(filename, "乐天") {
		res, err = Rakuten(content)
	} else if strings.Contains(filename, "Yahoo") {
		res, err = Yahoo(content)
	} else if strings.Contains(filename, "BAYCREW") {
		res, err = Baycrew(content)
	} else if strings.Contains(filename, "dot-st") {
		res, err = DotSt(content)
	} else if strings.Contains(filename, "nissen") {
		res, err = Nissen(content)
	} else if strings.Contains(filename, "SHOPLIST") {
		res, err = Shoplist(content)
	} else if strings.Contains(filename, "ZOZOTO") {
		res, err = Zozoto(content)
	} else {
		log.Fatal("不支持该文件")
	}

	if err != nil {
		log.Fatal(err)
	}

	outputName := strings.Replace(filepath.Base(filename), ext, ".json", 1)
	err = output(res, outputName)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(outputName + " 生成成功")
}

func output(res interface{}, filename string) error {
	b, _ := json.Marshal(res)
	return ioutil.WriteFile(filename, b, 0777)
}
