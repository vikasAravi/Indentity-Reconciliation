package reconciliation_service

import (
	"bitespeed/identity-reconciliation/schema"
	util "bitespeed/identity-reconciliation/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		payload, _ := json.Marshal(schema.PingResponse{
			Success: "pong",
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}

func GetContactDetails() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// API CONTEXT DETAILS
		logContext := util.BuildContext("Get Contact Details API")
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logContext)

		// DECODING THE REQUEST BODY
		decoder := json.NewDecoder(r.Body)
		var req schema.IdentityRequest
		err := decoder.Decode(&req)
		if err != nil {
			util.Log.WithFields(logContext).Error(fmt.Sprintf("identity request decode failure: %v",
				err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// SERVICE LAYER INVOKATION
		response, err := GetIdentityResponse(ctx, req)
		if err != nil {
			util.Log.WithFields(logContext).Error(fmt.Sprintf("error while fetching the contact details: %v",
				err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(response)
		w.Write(payload)
		return

	}

}
