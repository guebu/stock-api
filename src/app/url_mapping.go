package app

import (
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/stock-api/controllers"
	"net/http"
)

func mapUrls(){
	logger.Info("About to map URLs...", "app:Stock-Api", "Layer:App", "Func:StartApplication", "Status:Start")

	router.HandleFunc("/stock-items", controllers.ItemsController.Create).Methods(http.MethodPost)
	router.HandleFunc("/stock-items/{id}", controllers.ItemsController.Get).Methods(http.MethodGet)
	router.HandleFunc("/stock-items/search", controllers.ItemsController.Search).Methods(http.MethodPost)

	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
	router.HandleFunc("/pong", controllers.PingController.Pong).Methods(http.MethodGet)

	logger.Info("URLs mapped!", "app:Stock-Api", "Layer:App", "Func:StartApplication", "Status:End")
}
