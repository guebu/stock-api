package stock_items

import "time"

type StockItem struct {
	Id                string      `json:"id"`
	Type              string      `json:"type"`
	Manufacturer      string      `json:"manufacturer"`
	Label             string      `json:"label"`
	CreationDate      time.Time   `json:"created"`
	OnStockSince      time.Time   `json:"on_stock_since"`
	Description       Description `json:"description"`
	Note              string      `json:"note"`
	Picture           []Picture   `json:"pictures"`
	Price             float32     `json:"price"`
	AvailableQuantity int         `json:"available_quantity"`
	UsedQuantity      int         `json:"used_quantity"`
	Status            string      `json:"status"`
}

type Description struct {
	PlainText 	string `json:"plain_text"`
	Html 		string `json:"html"`
}

type Picture struct {
	Id 	int64 	`json:"id"`
	Url string 	`json:"url"`
}
