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

	challan, statusCode, errResp := FetchChallans(payload)
	if errResp.Error != "" {
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	json.NewEncoder(w).Encode(challan)
	w.WriteHeader(statusCode)
}

func FetchChallans(payload []byte) (challan *models.ChallanResponse, statusCode int, errResp *models.ErrorResponse) {

	var requestedURL string
	if os.Getenv("IN_PROD") == "1" {
		requestedURL = os.Getenv("PROD_URL") + os.Getenv("V1_CHALLAN_ENDPOINT")
	} else {
		requestedURL = os.Getenv("UAT_URL") + os.Getenv("V1_CHALLAN_ENDPOINT")
	}

	req, err := http.NewRequest("POST", requestedURL, bytes.NewBuffer(payload))
	if err != nil {
		errResp := models.NewErrorResponse("internal server error try after few hours")
		return &models.ChallanResponse{}, http.StatusInternalServerError, errResp
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Referer", "docs.apiclub.in")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errResp := models.NewErrorResponse("OOPS! something went wrong")
		return &models.ChallanResponse{}, http.StatusBadRequest, errResp
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	// Save body to file
	err = os.WriteFile("vehicle_challans.json", body, 0644)
	if err != nil {
		log.Println("Failed to save response:", err)
	}

	challan = models.NewVehicleChallanResponse()

	if err = json.Unmarshal(body, &challan); err != nil {
		errResp := models.NewErrorResponse("error in request")
		return &models.ChallanResponse{}, http.StatusExpectationFailed, errResp
	}

	switch res.StatusCode {
	case http.StatusOK:

		return challan, http.StatusOK, nil
	case http.StatusTooManyRequests:
		errResp := models.NewErrorResponse("request quota exceeded.")
		return &models.ChallanResponse{}, http.StatusTooManyRequests, errResp

	default:
		errResp := models.NewErrorResponse(fmt.Sprintf("message-> %v, code-> %d", challan.Message, challan.Code))
		return &models.ChallanResponse{}, res.StatusCode, errResp
	}
}
