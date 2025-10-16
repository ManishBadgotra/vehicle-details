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

func GetVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()
	var v models.VehicleDetails

	licensePlate := r.URL.Query().Get("license")
	vehicle, err := v.GetFromDB(licensePlate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(vehicle)
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()

	var (
		newVehicle models.VehicleDetails
		statusCode int
	)

	licensePlate := r.URL.Query().Get("license")
	chassis := r.URL.Query().Get("chassis")
	engine := r.URL.Query().Get("engine")

	if licensePlate == "" && chassis == "" && engine == "" {
		errResp := models.NewErrorResponse("vehicle, chassis and engine details are required")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	newVehicle, statusCode, err := FetchRcDetails(licensePlate, chassis, engine)
	if err != nil {
		errResponse := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(errResponse)
	}

	switch statusCode {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newVehicle)
	default:
		w.WriteHeader(statusCode)
	}

	// save returned result to db
	err = newVehicle.AddToDB()
	if err != nil {
		errResp := models.NewErrorResponse("unable to add details to database")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	// return rc details & challans details
	json.NewEncoder(w).Encode(newVehicle)
}

func UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()
	// get existing vehicle or return if not exist

	// fetch new detials

	// update detials to existing

	// update to db

	// return status of success or failure

	w.WriteHeader(200)
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()

	licensePlate := r.URL.Query().Get("license")

	var v models.VehicleDetails

	if err := v.DeleteFromDB(licensePlate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := models.NewErrorResponse("something went wrong. TRY AGAIN!")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	w.WriteHeader(200)
}

func FetchRcDetails(licensePlate, chassis, engine string) (newVehicle models.VehicleDetails, statusCode int, err error) {

	if licensePlate == "" {
		return models.VehicleDetails{}, http.StatusExpectationFailed, fmt.Errorf("vehicle number not provided")
	}

	request := models.NewRequestBody(licensePlate, chassis, engine)

	payload, err := json.Marshal(request)
	if err != nil {
		return models.VehicleDetails{}, http.StatusInternalServerError, fmt.Errorf("unable to create response for vehicle number")
	}

	req, err := http.NewRequest("POST", "https://uat.apiclub.in/api/v1/rc_info", bytes.NewBuffer(payload))
	if err != nil {
		return models.VehicleDetails{}, http.StatusInternalServerError, fmt.Errorf("unable to make request to the server")
	}

	req.Header.Add("x-api-key", apiKey)
	// req.Header.Add("x-Request-id", "") // adding request id is optional
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.VehicleDetails{}, http.StatusBadRequest, fmt.Errorf("something went wrong in requesting to server")
	}

	switch res.StatusCode {
	case http.StatusOK:

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))

		// Save body to file
		err = os.WriteFile("add_vehicle.json", body, 0644)
		if err != nil {
			log.Println("Failed to save response:", err)
		}

		if err = json.Unmarshal(body, &newVehicle); err != nil {
			return models.VehicleDetails{}, http.StatusNoContent, fmt.Errorf("no content received from server")
		}

	case http.StatusTooManyRequests:
		return models.VehicleDetails{}, http.StatusTooManyRequests, fmt.Errorf("request quota exceeded")
	default:
		return models.VehicleDetails{}, res.StatusCode, fmt.Errorf("error occured on third party request: %d", res.StatusCode)
	}

	// challanStruct, statusCode, errResp := FetchChallans(payload)
	// if errResp.Error != "" {
	// 	return models.VehicleDetails{}, statusCode, fmt.Errorf("%s", errResp.Error)
	// }

	// newVehicle.Response.Challans = challanStruct.Response

	return newVehicle, http.StatusOK, nil
}
