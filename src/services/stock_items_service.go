package services

import (
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/stock-api/domain/queries"
	"github.com/guebu/stock-api/domain/stock_items"
)

type stockItemServiceInterface interface {
	Create(stock_items.StockItem) (*stock_items.StockItem, *errors.ApplicationError)
	GetStockItemByID(string) (*stock_items.StockItem, *errors.ApplicationError)
	SearchStockItems(*queries.EsQuery) ([]stock_items.StockItem, *errors.ApplicationError)
}

type stockItemService struct {

}

var (
	ItemService stockItemServiceInterface = &stockItemService{}
)

func (sis *stockItemService) Create(stockItem stock_items.StockItem) (*stock_items.StockItem, *errors.ApplicationError) {

	err := stockItem.Save()

	if err != nil {
		return nil, err
	}
	return &stockItem, nil

}

func (sis *stockItemService) GetStockItemByID(id string) (*stock_items.StockItem, *errors.ApplicationError) {
	stockItem := stock_items.StockItem{Id: id}
	if err := stockItem.Get(); err != nil {
		return nil, err
	}
	return &stockItem, nil
}

func (sis *stockItemService) SearchStockItems(query *queries.EsQuery) ([]stock_items.StockItem, *errors.ApplicationError){
	dao := stock_items.StockItem{}
	return dao.Search(*query)

}

