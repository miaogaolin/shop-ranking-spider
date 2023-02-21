package main

type ProductPage struct {
	UpdateTime string        `json:"update_time"`
	List       []ProductList `json:"list"`
}
type ProductList struct {
	Site          string  `json:"-"`
	Key           string  `json:"-"`
	ID            string  `json:"id"`
	Url           string  `json:"url"`
	Ranking       uint64  `json:"ranking"`
	Image         string  `json:"image"`
	Title         string  `json:"title"`
	PriceBody     float64 `json:"price"`
	PriceYen      string  `json:"-"`
	StoreUrl      string  `json:"store_url"`
	StoreName     string  `json:"store_name"`
	StoreId       string  `json:"-"`
	SaleCount     uint64  `json:"-"`
	ReviewCount   uint64  `json:"review_count"`
	ReviewRanking float64 `json:"review_ranking"` // 评价星级
}
