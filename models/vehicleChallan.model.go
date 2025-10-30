package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/manishbadgotra/vehicle-details/database"
)

type ChallanResponse struct {
	Code     int      `json:"code"`
	Status   string   `json:"status"`
	Message  string   `json:"message"`
	Response Response `json:"response"`
}

type Response struct {
	RequestID   string    `json:"request_id"`
	VehicleID   string    `json:"vehicleId"`
	Total       int       `json:"total"`
	ChallanList []Challan `json:"challans"`
}

type Challan struct {
	ChallanNo     string `json:"challan_no"`
	Date          string `json:"date"`
	AccusedName   string `json:"accused_name"`
	ChallanStatus string `json:"challan_status"`
	Amount        int    `json:"amount"`
	State         string `json:"state"`
	Area          string `json:"area"`
	Offence       string `json:"offence"`
}

func NewVehicleChallanResponse() *ChallanResponse {
	return &ChallanResponse{}
}

func (challan ChallanResponse) Get() (int, int, []Challan, error) {
	var (
		total    int
		pending  int
		err      error
		c        Challan
		challans []Challan
	)

	conn, err := database.DBInstance.Conn(context.TODO())
	if err != nil {
		fmt.Fprintln(os.Stdout, "unable to establish connection")
		return total, pending, challans, fmt.Errorf("unable to establish connection")
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		fmt.Fprintln(os.Stdout, "unable to begin transaction")
		return total, pending, challans, fmt.Errorf("unable to begin transaction")
	}

	defer tx.Rollback()

	rows, err := tx.Query(`
			SELECT 
			challan_no, 
			date, 
			accused_name, 
			challan_status, 
			amount, 
			state, 
			area
			FROM challans
			WHERE license_plate = ?
		`,
		challan.Response.VehicleID,
	)
	if err != nil {
		return total, pending, challans, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&c.ChallanNo,
			&c.Date,
			&c.AccusedName,
			&c.ChallanStatus,
			&c.Amount,
			&c.State,
			&c.Area,
		); err != nil {
			return total, pending, challans, err
		}

		if c.ChallanStatus == "pending" || c.ChallanStatus == "Pending" {
			pending += 1
		}

		challans = append(challans, c)
	}

	total = len(challans)

	return total, pending, challans, err
}

func (challan ChallanResponse) Save() error {

	tx, err := database.DBInstance.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to establish connection")
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(
		database.ChallanInsert,
	)

	if err != nil {
		return fmt.Errorf("unable to prepare query")
	}

	defer stmt.Close()

	for _, c := range challan.Response.ChallanList {
		if _, err := stmt.Exec(
			&c.ChallanNo,
			&challan.Response.VehicleID,
			&c.Date,
			&c.AccusedName,
			&c.ChallanStatus,
			&c.Amount,
			&c.State,
			&c.Area,
			&c.Offence,
		); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (challan ChallanResponse) Update() error {

	if err := challan.Delete(); err != nil {
		return fmt.Errorf("unable to delete existing challans from database")
	}

	if err := challan.Save(); err != nil {
		return err
	}

	return nil
}

func (challan ChallanResponse) Delete() error {

	tx, err := database.DBInstance.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// first delete from `challans` table due to Foreign Key Constraints
	stmt, err := tx.Prepare(`
	DELETE FROM challans WHERE license_plate = ?
	`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(challan.Response.VehicleID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
