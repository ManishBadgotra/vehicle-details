package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/models"
)

func GetVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()

	licensePlate := r.URL.Query().Get("license")

	// open db connection
	db, err := database.OpenDB()
	if err != nil {
		w.WriteHeader(http.StatusMovedPermanently)
		fmt.Println(err)
		return
	}
	defer db.Close()

	// check for vehicle number
	rows, err := db.Query(
		`
		SELECT id, license_plate, owner_name, father_name, is_financed, financer, present_address, permanent_address,
		insurance_company, insurance_policy, insurance_expiry, class, registration_date, vehicle_age, pucc_upto, pucc_number,
		chassis_number, engine_number, fuel_type, brand_name, brand_model, cubic_capacity, gross_weight, cylinders, color, norms,
		noc_details, seating_capacity, owner_count, tax_upto, tax_paid_upto, permit_number, permit_issue_date, permit_valid_from,
		permit_valid_upto, permit_type, national_permit_number, national_permit_upto, national_permit_issued_by, rc_status FROM vehicles WHERE licenseplate == ?
		`,
		licensePlate,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	var v models.VehicleDetails

	var id int
	// get details of vehicle from db
	for rows.Next() {

		rows.Scan(&id, &v.Response.LicensePlate, v.Response.OwnerName, v.Response.FatherName, v.Response.IsFinanced, v.Response.Financer, v.Response.PresentAddress, v.Response.PermanentAddress,
			v.Response.InsuranceCompany, v.Response.InsurancePolicy, v.Response.InsuranceExpiry, v.Response.Class, v.Response.RegistrationDate, v.Response.VehicleAge, v.Response.PuccUpto, v.Response.PuccNumber,
			v.Response.ChassisNumber, v.Response.EngineNumber, v.Response.FuelType, v.Response.BrandName, v.Response.BrandModel, v.Response.CubicCapacity, v.Response.GrossWeight, v.Response.Cylinders, v.Response.Color, v.Response.Norms,
			v.Response.NocDetails, v.Response.SeatingCapacity, v.Response.OwnerCount, v.Response.TaxUpto, v.Response.TaxPaidUpto, v.Response.PermitNumber, v.Response.PermitIssueDate, v.Response.PermitValidFrom,
			v.Response.PermitValidUpto, v.Response.PermitType, v.Response.NationalPermitNumber, v.Response.PermitValidUpto, v.Response.NationalPermitIssuedBy, v.Response.RcStatus,
		)
	}

	json.NewEncoder(w).Encode(v)
	w.WriteHeader(200)
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
		json.NewEncoder(w).Encode(errResp)
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
		json.NewEncoder(w).Encode(errResp)
	}

	// return rc details & challans details
	json.NewEncoder(w).Encode(newVehicle)
}

func DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			str := fmt.Sprintf("- error in adding vehicle to Database with error: %v", recovered)
			log.Printf("%s", str)
		}
	}()

	// get vehicle number from db

	// get challans of that vehicle from db

	// delete all from db at once

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

	challanStruct, statusCode, errResp := FetchChallans(payload)
	if errResp.Error != "" {
		return models.VehicleDetails{}, statusCode, fmt.Errorf("%s", errResp.Error)
	}

	newVehicle.Response.Challans = challanStruct.Response

	return newVehicle, http.StatusOK, nil
}
