package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/manishbadgotra/vehicle-details/database"
)

type VehicleDetails struct {
	Code      int            `json:"code,omitempty"`
	Status    string         `json:"status,omitempty"`
	Message   string         `json:"message,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Response  rcVerification `json:"response,omitempty"`
}
type rcVerification struct {
	RequestID              string `json:"request_id,omitempty"`
	LicensePlate           string `json:"license_plate,omitempty"`
	OwnerName              string `json:"owner_name,omitempty"`
	FatherName             string `json:"father_name,omitempty"`
	IsFinanced             string `json:"is_financed,omitempty"`
	Financer               string `json:"financer,omitempty"`
	PresentAddress         string `json:"present_address,omitempty"`
	PermanentAddress       string `json:"permanent_address,omitempty"`
	InsuranceCompany       string `json:"insurance_company,omitempty"`
	InsurancePolicy        string `json:"insurance_policy,omitempty"`
	InsuranceExpiry        string `json:"insurance_expiry,omitempty"`
	Class                  string `json:"class,omitempty"`
	RegistrationDate       string `json:"registration_date,omitempty"`
	VehicleAge             any    `json:"vehicle_age,omitempty"`
	PuccUpto               string `json:"pucc_upto,omitempty"`
	PuccNumber             string `json:"pucc_number,omitempty"`
	ChassisNumber          string `json:"chassis_number,omitempty"`
	EngineNumber           string `json:"engine_number,omitempty"`
	FuelType               string `json:"fuel_type,omitempty"`
	BrandName              string `json:"brand_name,omitempty"`
	BrandModel             string `json:"brand_model,omitempty"`
	CubicCapacity          string `json:"cubic_capacity,omitempty"`
	GrossWeight            string `json:"gross_weight,omitempty"`
	Cylinders              string `json:"cylinders,omitempty"`
	Color                  string `json:"color,omitempty"`
	Norms                  string `json:"norms,omitempty"`
	NocDetails             string `json:"noc_details,omitempty"`
	SeatingCapacity        string `json:"seating_capacity,omitempty"`
	OwnerCount             string `json:"owner_count,omitempty"`
	TaxUpto                string `json:"tax_upto,omitempty"`
	TaxPaidUpto            string `json:"tax_paid_upto,omitempty"`
	PermitNumber           string `json:"permit_number,omitempty"`
	PermitIssueDate        string `json:"permit_issue_date,omitempty"`
	PermitValidFrom        string `json:"permit_valid_from,omitempty"`
	PermitValidUpto        string `json:"permit_valid_upto,omitempty"`
	PermitType             string `json:"permit_type,omitempty"`
	NationalPermitNumber   string `json:"national_permit_number,omitempty"`
	NationalPermitUpto     string `json:"national_permit_upto,omitempty"`
	NationalPermitIssuedBy string `json:"national_permit_issued_by,omitempty"`
	RcStatus               string `json:"rc_status,omitempty"`
	// Challan                VehicleChallans `json:"vehicle_challans"`
}

func (v *VehicleDetails) GetFromDB(licensePlate string) (VehicleDetails, error) {

	// var vehicles VehicleDetails
	conn, err := database.DBInstance.Conn(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stdout, "error 1")
		return VehicleDetails{}, fmt.Errorf("unable to establish connection")
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		fmt.Fprintln(os.Stdout, "error 2")
		return VehicleDetails{}, fmt.Errorf("unable to begin transaction")
	}

	// check in `vehicles` Table
	row := tx.QueryRow(
		database.FindInVehicleTable,
		licensePlate,
	)

	err = row.Scan(&v.Response.LicensePlate, &v.Response.OwnerName, &v.Response.FatherName, &v.Response.IsFinanced, &v.Response.Financer, &v.Response.PresentAddress, &v.Response.PermanentAddress,
		&v.Response.InsuranceCompany, &v.Response.InsurancePolicy, &v.Response.InsuranceExpiry, &v.Response.Class, &v.Response.RegistrationDate, &v.Response.VehicleAge, &v.Response.PuccUpto, &v.Response.PuccNumber,
		&v.Response.ChassisNumber, &v.Response.EngineNumber, &v.Response.FuelType, &v.Response.BrandName, &v.Response.BrandModel, &v.Response.CubicCapacity, &v.Response.GrossWeight, &v.Response.Cylinders, &v.Response.Color, &v.Response.Norms,
		&v.Response.NocDetails, &v.Response.SeatingCapacity, &v.Response.OwnerCount, &v.Response.TaxUpto, &v.Response.TaxPaidUpto, &v.Response.PermitNumber, &v.Response.PermitIssueDate, &v.Response.PermitValidFrom,
		&v.Response.PermitValidUpto, &v.Response.PermitType, &v.Response.NationalPermitNumber, &v.Response.PermitValidUpto, &v.Response.NationalPermitIssuedBy, &v.Response.RcStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(os.Stdout, "error 3")
			tx.Rollback()
			return VehicleDetails{}, fmt.Errorf("no record found")
		}
		// else {
		// 	fmt.Fprintln(os.Stdout, "error 4")
		// 	return VehicleDetails{}, err
		// }
	}

	return *v, nil
}

func (v *VehicleDetails) AddToDB() (err error) {
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

	if _, err := stmt.Exec(&v.Response.RequestID, &v.Response.LicensePlate, &v.Response.OwnerName, &v.Response.FatherName, &v.Response.IsFinanced, &v.Response.Financer, &v.Response.PresentAddress,
		&v.Response.InsuranceCompany, &v.Response.InsurancePolicy, &v.Response.InsuranceExpiry, &v.Response.Class, &v.Response.RegistrationDate, &v.Response.VehicleAge, &v.Response.PuccUpto, &v.Response.PuccNumber,
		&v.Response.ChassisNumber, &v.Response.EngineNumber, &v.Response.FuelType, &v.Response.BrandName, &v.Response.BrandModel, &v.Response.CubicCapacity, &v.Response.GrossWeight, &v.Response.Cylinders, &v.Response.Color, &v.Response.Norms,
		&v.Response.NocDetails, &v.Response.SeatingCapacity, &v.Response.OwnerCount, &v.Response.TaxUpto, &v.Response.TaxPaidUpto, &v.Response.PermitNumber, &v.Response.PermitIssueDate, &v.Response.PermitValidFrom,
		&v.Response.PermitValidUpto, &v.Response.PermitType, &v.Response.NationalPermitNumber, &v.Response.NationalPermitUpto, &v.Response.NationalPermitIssuedBy, &v.Response.RcStatus,
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

func (v *VehicleDetails) UpdateToDB() (err error) {

	// using transaction - Commit/Rollback pattern

	// get vehicles detial from `vehicles` table

	// get challans list from `challans` table

	// update all updated field appropraitly

	return nil
}

func (v *VehicleDetails) DeleteFromDB(licensePlate string) (err error) {

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
