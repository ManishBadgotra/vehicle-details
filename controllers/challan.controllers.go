package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/manishbadgotra/vehicle-details/models"
)

func GetVehicleChallans(w http.ResponseWriter, r *http.Request) {

	secretCode := r.Header.Get("x-request-code")
	if secretCode != "ManishIsAGenius" {
		errResponse := models.NewErrorResponse("secret key not setup correctly")
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	licensePlate := r.URL.Query().Get("license_plate")
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

	challanStruct, statusCode, errResp := FetchChallans(payload)
	if errResp.Error != "" {
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	switch statusCode {
	case http.StatusOK:
		json.NewEncoder(w).Encode(challanStruct)
		w.WriteHeader(statusCode)
	default:
		w.WriteHeader(statusCode)
	}
}

func FetchChallans(payload []byte) (challanStruct *models.VehicleChallans, statusCode int, errResp *models.ErrorResponse) {

	req, err := http.NewRequest("POST", "https://uat.apiclub.in/api/v1/rc_info", bytes.NewBuffer(payload))
	if err != nil {
		errResp := models.NewErrorResponse(err.Error())
		return &models.VehicleChallans{}, http.StatusInternalServerError, errResp
	}

	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errResp := models.NewErrorResponse(err.Error())
		return &models.VehicleChallans{}, http.StatusBadRequest, errResp
	}

	switch res.StatusCode {
	case http.StatusOK:
		challanStruct = models.NewVehicleChallanResponse()

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))

		// Save body to file
		err = os.WriteFile("vehicle_challans.json", body, 0644)
		if err != nil {
			log.Println("Failed to save response:", err)
		}

		if err = json.Unmarshal(body, &challanStruct); err != nil {
			errResp := models.NewErrorResponse(err.Error())
			return &models.VehicleChallans{}, http.StatusNoContent, errResp
		}

	case http.StatusTooManyRequests:
		errResp := models.NewErrorResponse("request quota exceeded.")
		return &models.VehicleChallans{}, http.StatusTooManyRequests, errResp

	default:
		errResp := models.NewErrorResponse(fmt.Sprintf("error occured on third party request: %d", res.StatusCode))
		return &models.VehicleChallans{}, res.StatusCode, errResp
	}

	return challanStruct, http.StatusOK, nil

}
