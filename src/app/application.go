package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/guebu/common-utils/logger"
	"github.com/guebu/stock-api/clients/elasticsearch"
	"net/http"
	"time"
)

var (
	router = mux.NewRouter()
)

func StartApplication() {
	logger.Info("About to start the application...", "app:Stock-Api", "Layer:App", "Func:StartApplication", "Status:Start")

	elasticsearch.Init()
	mapUrls()

	logger.Info("About to start the http server...", "app:Stock-Api", "Layer:App", "Func:StartApplication", "Status:Start")
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:8083",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Starting server...")
		panic(err)
	}


	logger.Info("Application started...", "app:Stock-Api", "Layer:App", "Func:StartApplication", "Status:End")
	return
}
