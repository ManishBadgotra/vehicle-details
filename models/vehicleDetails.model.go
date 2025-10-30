package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/manishbadgotra/vehicle-details/database"
)

type VehicleRequest struct {
	Code      int             `json:"code,omitempty"`
	Status    string          `json:"status,omitempty"`
	Message   string          `json:"message,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	Response  VehicleResponse `json:"response,omitempty"`
}
type VehicleResponse struct {
	RequestID        string `json:"request_id,omitempty"`
	LicensePlate     string `json:"license_plate,omitempty"`
	OwnerName        string `json:"owner_name,omitempty"`
	FatherName       string `json:"father_name,omitempty"`
	IsFinanced       bool   `json:"is_financed,omitempty"`
	Financer         string `json:"financer,omitempty"`
	PresentAddress   string `json:"present_address,omitempty"`
	PermanentAddress string `json:"permanent_address,omitempty"`
	InsuranceCompany string `json:"insurance_company,omitempty"`
	InsurancePolicy  string `json:"insurance_policy,omitempty"`
	InsuranceExpiry  string `json:"insurance_expiry,omitempty"`
	Class            string `json:"class,omitempty"`
	RegistrationDate string `json:"registration_date,omitempty"`
	// VehicleAge             *string   `json:"vehicle_age,omitempty"`
	PuccUpto      string `json:"pucc_upto,omitempty"`
	PuccNumber    string `json:"pucc_number,omitempty"`
	ChassisNumber string `json:"chassis_number,omitempty"`
	EngineNumber  string `json:"engine_number,omitempty"`
	FuelType      string `json:"fuel_type,omitempty"`
	BrandName     string `json:"brand_name,omitempty"`
	BrandModel    string `json:"brand_model,omitempty"`
	CubicCapacity string `json:"cubic_capacity,omitempty"`
	GrossWeight   string `json:"gross_weight,omitempty"`
	Cylinders     string `json:"cylinders,omitempty"`
	Color         string `json:"color,omitempty"`
	Norms         string `json:"norms,omitempty"`
	// NocDetails             string    `json:"noc_details,omitempty"`
	SeatingCapacity string `json:"seating_capacity,omitempty"`
	OwnerCount      string `json:"owner_count,omitempty"`
	Fitness         string `json:"fit_up_to,omitempty"`
	TaxUpto         string `json:"tax_upto,omitempty"`
	// TaxPaidUpto     string `json:"tax_paid_upto,omitempty"`
	PermitNumber string `json:"permit_number,omitempty"`
	// PermitIssueDate        *string   `json:"permit_issue_date,omitempty"`
	// PermitValidFrom        *string   `json:"permit_valid_from,omitempty"`
	PermitValidUpto        string    `json:"permit_valid_upto,omitempty"`
	PermitType             string    `json:"permit_type,omitempty"`
	NationalPermitNumber   *string   `json:"national_permit_number,omitempty"`
	NationalPermitUpto     *string   `json:"national_permit_upto,omitempty"`
	NationalPermitIssuedBy *string   `json:"national_permit_issued_by,omitempty"`
	TotalChallans          int       `json:"total_challans,omitempty"`
	PendingChallans        int       `json:"pending_challans,omitempty"`
	RcStatus               string    `json:"rc_status,omitempty"`
	ChallanList            []Challan `json:"challans,omitempty"`
}

func (v VehicleRequest) GetFromDB(licensePlate string) (VehicleResponse, error) {
	var vehicle VehicleResponse
	// var vehicles VehicleRequest
	conn, err := database.DBInstance.Conn(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stdout, "unable to establish connection")
		return vehicle, fmt.Errorf("unable to establish connection")
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		fmt.Fprintln(os.Stdout, "unable to begin transaction")
		return vehicle, fmt.Errorf("unable to begin transaction")
	}

	defer tx.Rollback()

	// check in `vehicles` Table
	row := tx.QueryRow(
		database.FindInVehicleTable,
		licensePlate,
	)

	err = row.Scan(
		&vehicle.LicensePlate,
		&vehicle.OwnerName,
		&vehicle.FatherName,
		&vehicle.IsFinanced,
		&vehicle.Financer,
		&vehicle.PresentAddress,
		&vehicle.PermanentAddress,
		&vehicle.InsuranceCompany,
		&vehicle.InsurancePolicy,
		&vehicle.InsuranceExpiry,
		&vehicle.Class,
		&vehicle.RegistrationDate,
		// &vehicle.VehicleAge,
		&vehicle.PuccUpto,
		&vehicle.PuccNumber,
		&vehicle.ChassisNumber,
		&vehicle.EngineNumber,
		&vehicle.FuelType,
		&vehicle.BrandName,
		&vehicle.BrandModel,
		&vehicle.CubicCapacity,
		&vehicle.GrossWeight,
		&vehicle.Cylinders,
		&vehicle.Color,
		&vehicle.Norms,
		// &vehicle.NocDetails,
		&vehicle.SeatingCapacity,
		&vehicle.OwnerCount,
		&vehicle.Fitness,
		&vehicle.TaxUpto,
		// &vehicle.TaxPaidUpto,
		&vehicle.PermitNumber,
		// &vehicle.PermitIssueDate,
		// &vehicle.PermitValidFrom,
		&vehicle.PermitValidUpto,
		&vehicle.PermitType,
		&vehicle.NationalPermitNumber,
		&vehicle.NationalPermitUpto,
		&vehicle.NationalPermitIssuedBy,
		&v.Response.TotalChallans,
		&v.Response.PendingChallans,
		&vehicle.RcStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(os.Stdout, "no record found in here")
			return vehicle, fmt.Errorf("no record found")
		}
		// else {
		// 	fmt.Fprintln(os.Stdout, "error 4")
		// 	return vehicle, err
		// }
	}

	if err = tx.Commit(); err != nil {
		return vehicle, fmt.Errorf("something went wrong")
	}

	return vehicle, nil
}

func (v *VehicleRequest) AddToDB() (err error) {
	tx, err := database.DBInstance.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to establish connection")
	}

	stmt, err := tx.Prepare(
		database.VehicleInsert,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to prepare query")
	}

	if _, err := stmt.Exec(
		&v.Response.LicensePlate,
		&v.Response.OwnerName,
		&v.Response.FatherName,
		&v.Response.IsFinanced,
		&v.Response.Financer,
		&v.Response.PresentAddress,
		&v.Response.PermanentAddress,
		&v.Response.InsuranceCompany,
		&v.Response.InsurancePolicy,
		&v.Response.InsuranceExpiry,
		&v.Response.Class,
		&v.Response.RegistrationDate,
		// &v.Response.VehicleAge,
		&v.Response.PuccUpto,
		&v.Response.PuccNumber,
		&v.Response.ChassisNumber,
		&v.Response.EngineNumber,
		&v.Response.FuelType,
		&v.Response.BrandName,
		&v.Response.BrandModel,
		&v.Response.CubicCapacity,
		&v.Response.GrossWeight,
		&v.Response.Cylinders,
		&v.Response.Color,
		&v.Response.Norms,
		// &v.Response.NocDetails,
		&v.Response.SeatingCapacity,
		&v.Response.OwnerCount,
		&v.Response.Fitness,
		&v.Response.TaxUpto,
		// &v.Response.TaxPaidUpto,
		&v.Response.PermitNumber,
		// &v.Response.PermitIssueDate,
		// &v.Response.PermitValidFrom,
		&v.Response.PermitValidUpto,
		&v.Response.PermitType,
		&v.Response.NationalPermitNumber,
		&v.Response.NationalPermitUpto,
		&v.Response.NationalPermitIssuedBy,
		&v.Response.TotalChallans,
		&v.Response.PendingChallans,
		&v.Response.RcStatus,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to execute request")
	}

	stmt.Close()

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v *VehicleRequest) UpdateToDB() error {

	// using transaction - Commit/Rollback pattern

	// get vehicles detial from `vehicles` table
	v.DeleteFromDB(v.Response.LicensePlate)

	err := v.AddToDB()
	if err != nil {
		log.Printf("error in adding vehicle number: %v to database, with error: %v", v.Response.LicensePlate, err)
		return err
	}

	return nil
}

func (v *VehicleRequest) DeleteFromDB(licensePlate string) (err error) {

	tx, err := database.DBInstance.Begin()
	if err != nil {
		return err
	}

	// first delete from `challans` table due to Foreign Key Constraints
	stmt, err := tx.Prepare(`
	DELETE FROM challans WHERE license_plate = ?
	`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err = stmt.Exec(licensePlate); err != nil {
		tx.Rollback()
		return err
	}

	// then from `vehicles` table
	stmt, err = tx.Prepare(`
	DELETE FROM vehicles WHERE license_plate = ?
	`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err = stmt.Exec(licensePlate); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func FetchVehicleDetails(payload []byte) (VehicleRequest, int, *ErrorResponse) {

	var (
		requestedURL string
		newVehicle   VehicleRequest
	)
	if os.Getenv("IN_PROD") == "1" {
		requestedURL = os.Getenv("PROD_URL") + os.Getenv("V1_VEHICLE_ENDPOINT")
	} else {
		requestedURL = os.Getenv("UAT_URL") + os.Getenv("V1_VEHICLE_ENDPOINT")
	}

	req, err := http.NewRequest("POST", requestedURL, bytes.NewBuffer(payload))
	if err != nil {
		errResp := NewErrorResponse("unable to make request to the server")
		return newVehicle, http.StatusInternalServerError, errResp
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Referer", "docs.apiclub.in")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errResp := NewErrorResponse(err.Error())
		return newVehicle, http.StatusBadRequest, errResp
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&newVehicle); err != nil {
		return newVehicle, http.StatusNotFound, NewErrorResponse(err.Error())
	}

	challanRequest := NewRequestBody(newVehicle.Response.LicensePlate, newVehicle.Response.ChassisNumber, newVehicle.Response.ChassisNumber)

	challanPayload, err := json.Marshal(challanRequest)
	if err != nil {
		// errResp := NewErrorResponse("unable to create response for vehicle number")
		return newVehicle, http.StatusInternalServerError, NewErrorResponse(err.Error())
	}

	challans, statusCode, errResp := FetchChallans(challanPayload)
	if errResp != nil && errResp.Error == "" {
		return newVehicle, statusCode, errResp
	}

	totalChallans, pendingChallans, challanLists, err := challans.Get()
	if err != nil {
		log.Printf("count challans error: %v", err)
	}

	if err == nil {
		newVehicle.Response.ChallanList = append(newVehicle.Response.ChallanList, challanLists...)
	}

	newVehicle.Response.TotalChallans = totalChallans
	newVehicle.Response.PendingChallans = pendingChallans

	err = newVehicle.UpdateToDB()
	if err != nil {
		errResp := NewErrorResponse("request successfull but unable to save data to database")
		return newVehicle, 503, errResp
	}

	return newVehicle, http.StatusOK, nil
}

func FetchRcDetails(licensePlate, chassis, engine string) (newVehicle VehicleRequest, statusCode int, err error) {

	if licensePlate == "" {
		return VehicleRequest{}, http.StatusExpectationFailed, fmt.Errorf("vehicle number not provided")
	}

	request := NewRequestBody(licensePlate, chassis, engine)

	payload, err := json.Marshal(request)
	if err != nil {
		return VehicleRequest{}, http.StatusInternalServerError, fmt.Errorf("unable to create response for vehicle number")
	}

	req, err := http.NewRequest("POST", "https://uat.apiclub.in/api/v1/rc_info", bytes.NewBuffer(payload))
	if err != nil {
		return VehicleRequest{}, http.StatusInternalServerError, fmt.Errorf("unable to make request to the server")
	}

	req.Header.Add("x-api-key", os.Getenv("API_KEY"))
	// req.Header.Add("x-Request-id", "") // adding request id is optional
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return VehicleRequest{}, http.StatusBadRequest, fmt.Errorf("something went wrong in requesting to server")
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
			return VehicleRequest{}, http.StatusNotFound, fmt.Errorf("no content received from server")
		}

	case http.StatusTooManyRequests:
		return VehicleRequest{}, http.StatusTooManyRequests, fmt.Errorf("request quota exceeded")
	default:
		return VehicleRequest{}, res.StatusCode, fmt.Errorf("error occured on third party request: %d", res.StatusCode)
	}

	// challanStruct, statusCode, errResp := FetchChallans(payload)
	// if errResp.Error != "" {
	// 	return VehicleRequest{}, statusCode, fmt.Errorf("%s", errResp.Error)
	// }

	// newVehicle.Response.Challans = challanStruct.Response

	return newVehicle, http.StatusOK, nil
}

func FetchChallans(payload []byte) (*ChallanResponse, int, *ErrorResponse) {

	var (
		errResp *ErrorResponse
		challan *ChallanResponse
	)

	var requestedURL string
	if os.Getenv("IN_PROD") == "1" {
		requestedURL = os.Getenv("PROD_URL") + os.Getenv("V1_CHALLAN_ENDPOINT")
	} else {
		requestedURL = os.Getenv("UAT_URL") + os.Getenv("V1_CHALLAN_ENDPOINT")
	}

	req, err := http.NewRequest("POST", requestedURL, bytes.NewBuffer(payload))
	if err != nil {
		errResp = NewErrorResponse("internal server error try after few hours")
		return nil, http.StatusInternalServerError, errResp
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Referer", "docs.apiclub.in")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errResp = NewErrorResponse("OOPS! something went wrong")
		return nil, http.StatusBadRequest, errResp
	}

	// fmt.Println(string(body))

	// // Save body to file
	// err = os.WriteFile("vehicle_challans.json", body, 0644)
	// if err != nil {
	// 	log.Println("Failed to save response:", err)
	// }

	challan = NewVehicleChallanResponse()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&challan)
	if err != nil {
		errResp = NewErrorResponse("error in request")
		return nil, http.StatusExpectationFailed, errResp
	}

	switch res.StatusCode {
	case http.StatusOK:

		// delete existing challans of same vehicle
		challan.Delete()
		// save new challans again in db
		challan.Save()

		return challan, http.StatusOK, nil
	case http.StatusTooManyRequests:
		errResp = NewErrorResponse("request quota exceeded.")
		return nil, http.StatusTooManyRequests, errResp

	default:
		errResp = NewErrorResponse(fmt.Sprintf("message-> %v, code-> %d", challan.Message, challan.Code))
		return nil, res.StatusCode, errResp
	}
}
