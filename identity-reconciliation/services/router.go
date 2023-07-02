package services

import (
	"bitespeed/identity-reconciliation/services/reconciliation_service"
	"github.com/gorilla/mux"

	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", reconciliation_service.PingHandler()).Methods("GET")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	groupV1 := router.PathPrefix("/api/v1").Subrouter()
	groupV1.HandleFunc("/identify", reconciliation_service.GetContactDetails()).Methods("POST")
	return router
}
