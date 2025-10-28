package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/manishbadgotra/vehicle-details/models"
)

func GetVehicle(w http.ResponseWriter, r *http.Request) {
	v := models.VehicleRequest{}

	licensePlate := r.URL.Query().Get("license")
	existingVehicle, err := v.GetFromDB(licensePlate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errResp := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(existingVehicle)
}

func AddVehicle(w http.ResponseWriter, r *http.Request) {
	var requestedURL string
	if os.Getenv("IN_PROD") == "1" {
		requestedURL = os.Getenv("PROD_URL") + os.Getenv("V1_VEHICLE_ENDPOINT")
	} else {
		requestedURL = os.Getenv("UAT_URL") + os.Getenv("V1_VEHICLE_ENDPOINT")
	}

	reqBody := vehicleStruct{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		errResp := models.NewErrorResponse("request unsupported")
		json.NewEncoder(w).Encode(errResp)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.String("Path Params --> ", reqBody.VehicleId)

	payload, err := json.Marshal(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := models.NewErrorResponse("unable to create response for vehicle number")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req, err := http.NewRequest("POST", requestedURL, bytes.NewBuffer(payload))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := models.NewErrorResponse("unable to make request to the server")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Referer", "docs.apiclub.in")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// defer res.Body.Close()

	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(models.NewErrorResponse("unable to read response body"))
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	decoder = json.NewDecoder(res.Body)
	// w.WriteHeader(res.StatusCode)
	// w.Write(body)

	vehicleStruct := models.VehicleRequest{}

	if err := decoder.Decode(&vehicleStruct); err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// // Save body to file
	// err = os.WriteFile("vehicle_details.json", body, 0644)
	// if err != nil {
	// 	log.Println("Failed to save response:", err)
	// }

	// if err = json.Unmarshal(body, &vehicleStruct); err != nil {
	// 	w.WriteHeader(http.StatusNoContent)
	// 	return
	// }

	err = vehicleStruct.AddToDB()
	if err != nil {
		w.WriteHeader(503)

		errResp := models.NewErrorResponse("request successfull but unable to save data to database")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("Vehicle Details Request successfull")
	json.NewEncoder(w).Encode(vehicleStruct.Response)
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

	var v models.VehicleRequest

	if err := v.DeleteFromDB(licensePlate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResp := models.NewErrorResponse("something went wrong. TRY AGAIN!")
		json.NewEncoder(w).Encode(&errResp)
		return
	}

	w.WriteHeader(200)
}

func FetchRcDetails(licensePlate, chassis, engine string) (newVehicle models.VehicleRequest, statusCode int, err error) {

	if licensePlate == "" {
		return models.VehicleRequest{}, http.StatusExpectationFailed, fmt.Errorf("vehicle number not provided")
	}

	request := models.NewRequestBody(licensePlate, chassis, engine)

	payload, err := json.Marshal(request)
	if err != nil {
		return models.VehicleRequest{}, http.StatusInternalServerError, fmt.Errorf("unable to create response for vehicle number")
	}

	req, err := http.NewRequest("POST", "https://uat.apiclub.in/api/v1/rc_info", bytes.NewBuffer(payload))
	if err != nil {
		return models.VehicleRequest{}, http.StatusInternalServerError, fmt.Errorf("unable to make request to the server")
	}

	req.Header.Add("x-api-key", apiKey)
	// req.Header.Add("x-Request-id", "") // adding request id is optional
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.VehicleRequest{}, http.StatusBadRequest, fmt.Errorf("something went wrong in requesting to server")
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
			return models.VehicleRequest{}, http.StatusNoContent, fmt.Errorf("no content received from server")
		}

	case http.StatusTooManyRequests:
		return models.VehicleRequest{}, http.StatusTooManyRequests, fmt.Errorf("request quota exceeded")
	default:
		return models.VehicleRequest{}, res.StatusCode, fmt.Errorf("error occured on third party request: %d", res.StatusCode)
	}

	// challanStruct, statusCode, errResp := FetchChallans(payload)
	// if errResp.Error != "" {
	// 	return models.VehicleRequest{}, statusCode, fmt.Errorf("%s", errResp.Error)
	// }

	// newVehicle.Response.Challans = challanStruct.Response

	return newVehicle, http.StatusOK, nil
}
