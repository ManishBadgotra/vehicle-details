package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/manishbadgotra/vehicle-details/models"
)

func GetVehicleChallans(w http.ResponseWriter, r *http.Request) {

	// secretCode := r.Header.Get("x-request-code")
	// if secretCode != "ManishIsAGenius" {
	// 	errResponse := models.NewErrorResponse("secret key not setup correctly")
	// 	json.NewEncoder(w).Encode(errResponse)
	// 	return
	// }

	licensePlate := r.URL.Query().Get("license")
	chassis := r.URL.Query().Get("chassis")
	engine := r.URL.Query().Get("engine")

	if licensePlate == "" && chassis == "" && engine == "" {
		errResp := models.NewErrorResponse("vehicle, chassis and engine details are required")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	request := models.NewRequestBody(licensePlate, chassis, engine)

	payload, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := models.NewErrorResponse("unable to create response for vehicle number")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	challan, statusCode, errResp := models.FetchChallans(payload)
	if errResp != nil && errResp.Error == "" {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(challan)
}
