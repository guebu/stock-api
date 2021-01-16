package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	errors2 "github.com/guebu/common-utils/errors"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/oauth-go/oauth"
	"github.com/guebu/stock-api/domain/queries"
	"github.com/guebu/stock-api/domain/stock_items"
	"github.com/guebu/stock-api/services"
	"github.com/guebu/stock-api/utils/http_utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Search(http.ResponseWriter, *http.Request)

}

type itemsController struct {

}

var (
	ItemsController = itemsController{}
)

func (itc *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("About to create a stock item...", "app:stock-api", "layer:controller", "func:Create", "status:start")

	stockItem, err := itc.getStockItemDataFromRequest(w, r)

	if err != nil {
		//Return it back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.AStatusCode)
		json.NewEncoder(w).Encode(err)
		logger.Error("Error while trying to read stock item information from request", err, "app:app:stock-api", "layer:controller", "func:Create", "status:error")
		return
	}

	// stock item data could be read successfully
	// Now chek if request was authorized
	err = oauth.AuthenticateRequest(r)

	if err != nil {
		if  err.AnError == "token already expired" {
			//ToDo: Handle that token is already expired
			logger.Error("Token already expired!", err, "app:stock-api", "layer:controller", "func:Create", "status:in progress")
			http_utils.RespondErrorAsJson(w, err)
			return
		}
		//Another error occured!
		//Return it back to the client
		http_utils.RespondErrorAsJson(w, err)
		return
	}

	//Request is authorized
	logger.Info("request is authorized...", "app:stock-api", "layer:controller", "func:Get", "status:in progress")
	http_utils.RespondBodyAsJson(w, http.StatusOK, stockItem)

	stockItem, err = services.ItemService.Create(*stockItem)
	if err != nil {
		http_utils.RespondErrorAsJson(w, err)
		return
	}
	http_utils.RespondBodyAsJson(w, http.StatusOK, stockItem)

	logger.Info("stock item created successfully...", "app:stock-api", "layer:controller", "func:Create", "status:end")
	return

}

func (itc *itemsController) Get(w http.ResponseWriter, r *http.Request) {

	logger.Info("About to start getting the stock item...", "app:stock-api", "layer:controller", "func:Get", "status:start")
	err := oauth.AuthenticateRequest(r)

	if err != nil {
		if err.AnError == "token already expired" {
			//ToDo: Handle that token is already expired

			logger.Info("Token expired!")
		}
		//An error occured!
		//Return it back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.AStatusCode)
		json.NewEncoder(w).Encode(err)
		return
	}

	//Request is authorized
	logger.Info("request is authorized...", "app:stock-api", "layer:controller", "func:Get", "status:in progress")

	vars := mux.Vars(r)
	itemID := strings.TrimSpace(vars["id"])

	stockItem, err := services.ItemService.GetStockItemByID(itemID)

	if err != nil {
		http_utils.RespondErrorAsJson(w,err)
		return
	}
	http_utils.RespondBodyAsJson(w, http.StatusOK, stockItem)

}

func (itc *itemsController) getStockItemDataFromRequest(w http.ResponseWriter, r *http.Request) (*stock_items.StockItem, *errors2.ApplicationError) {

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		err := errors2.NewBadRequestError("Not the right content type in request for creating stock item!", nil)
		//Return it back to the client
		return nil, err
	}

	var stockItem stock_items.StockItem
	stockItem = stock_items.StockItem{}

	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&stockItem)

	if err != nil {
		if errors.As(err, &unmarshalErr) {
			err := errors2.NewBadRequestError("Bad Request. Wrong Type provided for field " + unmarshalErr.Field, err)
			//Return it back to the client
			return nil, err
		}

		err := errors2.NewInternalServerError("Other error occured!", err)
		//Return it back to the client
		return nil, err
	}

	logger.Info("stock item created successfully...", "app:stock-api", "layer:controller", "func:Create", "status:end")
	return &stockItem, nil
}

func (c *itemsController) Search(w http.ResponseWriter, r *http.Request) {
	var query queries.EsQuery

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		apiErr := errors2.NewBadRequestError("Invalid JSON for stock items query!", nil )
		http_utils.RespondErrorAsJson(w, apiErr)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(bytes, &query); err != nil {
		apiErr := errors2.NewBadRequestError("Invalid JSON for stock items query!", nil )
		http_utils.RespondErrorAsJson(w, apiErr)
		return
	}

	stockItems, searchErr := services.ItemService.SearchStockItems(&query)

	if searchErr != nil {
		http_utils.RespondErrorAsJson(w, searchErr)
		return
	}

	http_utils.RespondBodyAsJson(w, http.StatusOK, stockItems)


}
