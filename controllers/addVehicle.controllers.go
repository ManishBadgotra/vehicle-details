package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manishbadgotra/vehicle-details/models"
	"github.com/manishbadgotra/vehicle-details/utils"
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

	reqBody := vehicleStruct{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqBody); err != nil {
		errResp := models.NewErrorResponse("request unsupported")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	// slog.String("Path Params --> ", reqBody.VehicleId)
	ok := utils.VerifyVehicleNumber(reqBody.VehicleId)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errResp := models.NewErrorResponse("enter valid vehicle number")

		json.NewEncoder(w).Encode(errResp)
		return
	}

	payload, err := json.Marshal(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := models.NewErrorResponse("unable to create response for vehicle number")
		json.NewEncoder(w).Encode(errResp)
		return
	}
	// v := models.VehicleRequest{}
	// existingVehicle, err := v.GetFromDB(reqBody.VehicleId)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	errResp := models.NewErrorResponse(err.Error())
	// 	json.NewEncoder(w).Encode(&errResp)
	// 	return
	// }

	// if existingVehicle.LicensePlate != "" {
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(existingVehicle)
	// 	return
	// }

	newVehicle, statusCode, errResp := models.FetchVehicleDetails(payload)
	if errResp != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newVehicle.Response)
	log.Println("Vehicle Details Request successfull")
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
