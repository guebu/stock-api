package controllers

import (
	"github.com/guebu/common-utils/controller"
	"net/http"
)

const (
	pingMessage = "Ping!"
	pongMessage = "Pong!"
)

type pingControllerInterface interface {
	Ping(http.ResponseWriter, *http.Request)
	Pong(http.ResponseWriter, *http.Request)
}

type pingController struct {

}

var (
	PingController = pingController{}
)

func (ping *pingController) Ping(w http.ResponseWriter, r *http.Request) {
	//Important: You have to set the content type first!! Otherise you will always get text/plain as content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(controller.GetJSONMessage(pingMessage)))
}


func (ping *pingController) Pong(w http.ResponseWriter, r *http.Request) {
	//Important: You have to set the content type first!! Otherise you will always get text/plain as content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(controller.GetJSONMessage(pongMessage)))

}
