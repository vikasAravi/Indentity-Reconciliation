package services

import (
	"bitespeed/identity-reconciliation/config"
	util "bitespeed/identity-reconciliation/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
	"time"
)

func StartAPI(c *cli.Context) error {
	router := Router()
	startServer(router)
	return nil
}

func startServer(router *mux.Router) {
	server := negroni.New(negroni.NewRecovery())
	handlerFunc := router.ServeHTTP
	server.UseHandlerFunc(handlerFunc)
	server.Use(httpStatLogger())
	portInfo := ":" + strconv.Itoa(int(config.Port()))
	util.Log.Info("Starting http server on port ", portInfo)
	e := http.ListenAndServe(portInfo, server)
	if e != nil {
		panic(e)
	}
	util.Log.Info("Stopping http server")
}

func httpStatLogger() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		startTime := time.Now()
		next(rw, r)
		responseTime := time.Now()
		deltaTime := responseTime.Sub(startTime).Seconds() * 1000

		if r.URL.Path != "/ping" {
			util.Log.WithFields(logrus.Fields{
				"RequestTime":   startTime.Format(time.RFC3339),
				"ResponseTime":  responseTime.Format(time.RFC3339),
				"DeltaTime":     deltaTime,
				"RequestURL":    r.URL.Path,
				"RequestMethod": r.Method,
			}).Debug("Http Logs")
		}
	}
}
