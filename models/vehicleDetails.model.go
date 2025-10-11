package models

import (
	"database/sql"
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
	Challans               ChallanResponse `json:"challans,omitempty"`
}

func (v *VehicleDetails) GetFromDB() error {

	

	return nil
}

func (v *VehicleDetails) AddToDB() (err error) {
	r := v.Response

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
func (v *VehicleDetails) UpdateFromDB() error {
	return nil
}
func (v *VehicleDetails) DeleteFromDB() error {
	return nil
}
