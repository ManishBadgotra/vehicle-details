package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/manishbadgotra/vehicle-details/database"
)

var db *sql.DB

type VehicleDetails struct {
	Code      int            `json:"code,omitempty"`
	Status    string         `json:"status,omitempty"`
	Message   string         `json:"message,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Response  rcVerification `json:"response,omitempty"`
}
type rcVerification struct {
	RequestID              string          `json:"request_id,omitempty"`
	LicensePlate           string          `json:"license_plate,omitempty"`
	OwnerName              string          `json:"owner_name,omitempty"`
	FatherName             string          `json:"father_name,omitempty"`
	IsFinanced             string          `json:"is_financed,omitempty"`
	Financer               string          `json:"financer,omitempty"`
	PresentAddress         string          `json:"present_address,omitempty"`
	PermanentAddress       string          `json:"permanent_address,omitempty"`
	InsuranceCompany       string          `json:"insurance_company,omitempty"`
	InsurancePolicy        string          `json:"insurance_policy,omitempty"`
	InsuranceExpiry        string          `json:"insurance_expiry,omitempty"`
	Class                  string          `json:"class,omitempty"`
	RegistrationDate       string          `json:"registration_date,omitempty"`
	VehicleAge             any             `json:"vehicle_age,omitempty"`
	PuccUpto               string          `json:"pucc_upto,omitempty"`
	PuccNumber             string          `json:"pucc_number,omitempty"`
	ChassisNumber          string          `json:"chassis_number,omitempty"`
	EngineNumber           string          `json:"engine_number,omitempty"`
	FuelType               string          `json:"fuel_type,omitempty"`
	BrandName              string          `json:"brand_name,omitempty"`
	BrandModel             string          `json:"brand_model,omitempty"`
	CubicCapacity          string          `json:"cubic_capacity,omitempty"`
	GrossWeight            string          `json:"gross_weight,omitempty"`
	Cylinders              string          `json:"cylinders,omitempty"`
	Color                  string          `json:"color,omitempty"`
	Norms                  string          `json:"norms,omitempty"`
	NocDetails             string          `json:"noc_details,omitempty"`
	SeatingCapacity        string          `json:"seating_capacity,omitempty"`
	OwnerCount             string          `json:"owner_count,omitempty"`
	TaxUpto                string          `json:"tax_upto,omitempty"`
	TaxPaidUpto            string          `json:"tax_paid_upto,omitempty"`
	PermitNumber           string          `json:"permit_number,omitempty"`
	PermitIssueDate        string          `json:"permit_issue_date,omitempty"`
	PermitValidFrom        string          `json:"permit_valid_from,omitempty"`
	PermitValidUpto        string          `json:"permit_valid_upto,omitempty"`
	PermitType             string          `json:"permit_type,omitempty"`
	NationalPermitNumber   string          `json:"national_permit_number,omitempty"`
	NationalPermitUpto     string          `json:"national_permit_upto,omitempty"`
	NationalPermitIssuedBy string          `json:"national_permit_issued_by,omitempty"`
	RcStatus               string          `json:"rc_status,omitempty"`
	Challans               ChallanResponse `json:"challans_response"`
}

func (v *VehicleDetails) GetFromDB(licensePlate, chassis, engine string) (VehicleDetails, error) {
	// open db connection
	db, err := database.OpenDB()
	if err != nil {
		return VehicleDetails{}, err
	}
	defer db.Close()

	// check in `vehicles` Table
	row := db.QueryRow(
		database.FindInVehicleTable,
		licensePlate,
		chassis,
		engine,
	)

	var id int

	err = row.Scan(&id, &v.Response.LicensePlate, v.Response.OwnerName, v.Response.FatherName, v.Response.IsFinanced, v.Response.Financer, v.Response.PresentAddress, v.Response.PermanentAddress,
		v.Response.InsuranceCompany, v.Response.InsurancePolicy, v.Response.InsuranceExpiry, v.Response.Class, v.Response.RegistrationDate, v.Response.VehicleAge, v.Response.PuccUpto, v.Response.PuccNumber,
		v.Response.ChassisNumber, v.Response.EngineNumber, v.Response.FuelType, v.Response.BrandName, v.Response.BrandModel, v.Response.CubicCapacity, v.Response.GrossWeight, v.Response.Cylinders, v.Response.Color, v.Response.Norms,
		v.Response.NocDetails, v.Response.SeatingCapacity, v.Response.OwnerCount, v.Response.TaxUpto, v.Response.TaxPaidUpto, v.Response.PermitNumber, v.Response.PermitIssueDate, v.Response.PermitValidFrom,
		v.Response.PermitValidUpto, v.Response.PermitType, v.Response.NationalPermitNumber, v.Response.PermitValidUpto, v.Response.NationalPermitIssuedBy, v.Response.RcStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return VehicleDetails{}, fmt.Errorf("no record found")
		} else {
			return VehicleDetails{}, err
		}
	}

	rows, err := db.Query(
		database.FindInChallansTable,
		v.Response.LicensePlate,
	)
	if err != nil {
		return *v, err
	}

	var challanId int

	for rows.Next() {
		challan := NewChallan()
		var offenceListJoined string
		rows.Scan(&challanId, &challan.ChallanNo, &challan.Date, &challan.AccusedName, &challan.ChallanStatus, &challan.Amount, &challan.State, &challan.Area, &challan.Offence, &offenceListJoined)
		offenceList := strings.Split(offenceListJoined, ",")
		for i, offenceName := range offenceList {
			challan.OffenceList[i].OffenceName = offenceName
			v.Response.Challans.Challans[i] = *challan
		}
	}

	return *v, nil
}

func (v *VehicleDetails) AddToDB() (err error) {
	r := &v.Response // this might create a copy of vehiceDetials's copy that can ignore updateing in original variable

	db, err = database.OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	vehicleStatement, err := tx.Prepare(database.VehicleInsert)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer vehicleStatement.Close()
	_, err = vehicleStatement.Exec(
		r.RequestID, r.LicensePlate, r.OwnerName, r.FatherName, r.IsFinanced, r.Financer, r.PresentAddress,
		r.InsuranceCompany, r.InsurancePolicy, r.InsuranceExpiry, r.Class, r.RegistrationDate, r.VehicleAge, r.PuccUpto, r.PuccNumber,
		r.ChassisNumber, r.EngineNumber, r.FuelType, r.BrandName, r.BrandModel, r.CubicCapacity, r.GrossWeight, r.Cylinders, r.Color, r.Norms,
		r.NocDetails, r.SeatingCapacity, r.OwnerCount, r.TaxUpto, r.TaxPaidUpto, r.PermitNumber, r.PermitIssueDate, r.PermitValidFrom,
		r.PermitValidUpto, r.PermitType, r.NationalPermitNumber, r.NationalPermitUpto, r.NationalPermitIssuedBy, r.RcStatus,
	)

	if err != nil {
		return err
	}

	challanStatement, err := tx.Prepare(database.ChallanInsert)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, challan := range r.Challans.Challans {
		offenceNames := []string{}
		for _, o := range challan.OffenceList {
			offenceNames = append(offenceNames, o.OffenceName)
		}
		offenceListStr := strings.Join(offenceNames, ",")

		challanStatement.Exec(
			challan.ChallanNo,
			r.LicensePlate,
			challan.Date,
			challan.AccusedName,
			challan.ChallanStatus,
			challan.Amount,
			challan.State,
			challan.Area,
			challan.Offence,
			offenceListStr,
		)
	}

	defer challanStatement.Close()

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (v *VehicleDetails) UpdateToDB() error {

	// using transaction - Commit/Rollback pattern

	// get vehicles detial from `vehicles` table

	// get challans list from `challans` table

	// update all updated field appropraitly

	return nil
}

func (v *VehicleDetails) DeleteFromDB(licensePlate string) error {

	db, err := database.OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
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
