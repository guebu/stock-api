package stock_items

import (
	"encoding/json"
	"fmt"
	"github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/stock-api/clients/elasticsearch"
	"github.com/guebu/stock-api/domain/queries"
	"strings"
)

//ToDo: Migrate the const to somewhere else... config file or whatever!!!
const (
	indexStockItems 	= "stockitems"
	docTypeStockItems 	= "stockitems"
)

func (si *StockItem) Save() *errors.ApplicationError {
	logger.Info("About to start saving the stock item", "app:StockItem API", "Layer:Domain", "Func:Save", "Status:Start")
	result, err := elasticsearch.Client.Index(indexStockItems, docTypeStockItems, si)

	if err != nil {
		logger.Error("Error while saving the StockItem in Elastic Search!", err, "app:StockItem API", "Layer:Domain", "Func:Save", "Status:Error")
		return errors.NewInternalServerError("Error while saving the StockItem in Elastic Search!", err)
	}
	si.Id = result.Id
	logger.Info(fmt.Sprintln("Stock Item saved successfully with ID %d", si.Id), "app:StockItem API", "Layer:Domain", "Func:Save", "Status:Start")
	return nil
}

func (si *StockItem) Get() *errors.ApplicationError {

	logger.Info(fmt.Sprintf("Start to get Stock item with given ID %s!", si.Id), "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Start")

	itemID := si.Id

	result, err := elasticsearch.Client.Get(indexStockItems, docTypeStockItems, si.Id)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			notFoundError := errors.NewNotFoundError(fmt.Sprintf("Stock item type with given ID %s wasn't found!", si.Id), nil )
			logger.Error(fmt.Sprintf("Stock item type with given ID %s wasn't found!", si.Id), notFoundError, "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Error")
			return notFoundError
		}
		logger.Error(fmt.Sprintf("Error while getting the stock item information with given ID %s", si.Id), err, "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Error" )
		return errors.NewInternalServerError(fmt.Sprintf("Error while getting the stock item information with given ID %s", si.Id), err)
	}

	if !result.Found {
		notFoundError := errors.NewNotFoundError(fmt.Sprintf("Stock item type with given ID %s wasn't found!", si.Id), nil )
		logger.Error(fmt.Sprintf("Stock item type with given ID %s wasn't found!", si.Id), notFoundError, "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Error")
		return notFoundError
	}

	bytes, marshallErr := result.Source.MarshalJSON()
	if err := json.Unmarshal(bytes, si); ( err != nil || marshallErr != nil ) {
		if err != nil {
			logger.Error(fmt.Sprintf("Error when trying to parse database response for stock items based on give ID %s ", si.Id), err, "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Error")
			return errors.NewInternalServerError(fmt.Sprintf("Error when trying to parse database response for stock items based on give ID %s ", si.Id), err)
		} else {
			logger.Error(fmt.Sprintf("Error when trying to parse database response for stock items based on give ID %s ", si.Id), marshallErr, "app:StockItem API", "Layer:Domain", "Func:Get", "Status:Error")
			return errors.NewInternalServerError(fmt.Sprintf("Error when trying to parse database response for stock items based on give ID %s ", si.Id), marshallErr)
		}
	}

	si.Id = itemID

	logger.Info(fmt.Sprintf("Stock item with given ID %s successfully found!", si.Id), "app:StockItem API", "Layer:Domain", "Func:Get", "Status:End")
	return nil
}

func (si *StockItem) Search(query queries.EsQuery) ([]StockItem, *errors.ApplicationError) {
	logger.Info("Start to search Stock items!", "app:stockitem api", "layer:domain", "func:search", "status:start")

	result, err := elasticsearch.Client.Search(indexStockItems, query.Build())

	if err != nil {
		logger.Error("Error when trying to search Stock item documents!", err, "app:stockitem api", "layer:domain", "func:search", "status:error")
		return nil, errors.NewInternalServerError("Error when trying to search Stock item documents!", err )
	}

	fmt.Println(result)

	stockItems := make([]StockItem, result.TotalHits())

	for i, hit := range result.Hits.Hits {
		bytes, err := hit.Source.MarshalJSON()
		if err != nil {
			logger.Error("Error during Marshaling JSON from search query!", err, "app:stockitem api", "layer:domain", "func:search", "status:error")
			return nil, errors.NewInternalServerError("Error during unmarshalling result from search query!", err)
		}
		var stockItem StockItem
		if err := json.Unmarshal(bytes, &stockItem); err != nil {
			logger.Error("Error during unmarshalling result in stock items from search query!", err, "app:StockItem API", "Layer:Domain", "Func:Search", "Status:Error")
			return nil, errors.NewInternalServerError("Error during unmarshalling result from search query!", err)
		}
		stockItems[i] = stockItem
	}

	if len(stockItems) == 0 {
		err := errors.NewNotFoundError("No matching stock items found!", nil )
		logger.Error("No matching stock items found!", err, "app:stockitem api", "layer:domain", "func:Search", "status:Error")
		return nil, err
	}

	logger.Info("Searched stock items successfully!", "app:stockitem api", "layer:domain", "func:search", "status:end")
	return stockItems, nil
}