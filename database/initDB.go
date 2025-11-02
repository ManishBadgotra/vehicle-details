package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // this package is required for sqlite driver
)

var (
	FindInVehicleTable = `
		SELECT 
		license_plate, 
		owner_name, 
		father_name, 
		is_financed, 
		financer, 
		present_address, 
		permanent_address,
		insurance_company, 
		insurance_policy, 
		insurance_expiry, 
		class, 
		registration_date,
		pucc_upto, 
		pucc_number,
		chassis_number, 
		engine_number, 
		fuel_type, 
		brand_name, 
		brand_model, 
		cubic_capacity, 
		gross_weight, 
		cylinders, 
		color, 
		norms,
		seating_capacity, 
		owner_count,
		fitness, 
		tax_upto, 
		permit_number, 
		permit_valid_upto, 
		permit_type, 
		national_permit_number, 
		national_permit_upto, 
		national_permit_issued_by, 
		total_challans, 
		pending_challans,
		rc_status
		FROM vehicles 
		WHERE license_plate = ?
	`
	FindInChallansTable = `
		SELECT 
		challan_no, 
		date, 
		accused_name, 
		challan_status, 
		amount, 
		state, 
		area, 
		offence 
		FROM challans WHERE license_plate = ?
	`
	VehicleInsert = `
        INSERT INTO vehicles
            (
			license_plate, 
			owner_name, 
			father_name, 
			is_financed, 
			financer, 
			present_address, 
			permanent_address,
            insurance_company, 
			insurance_policy, 
			insurance_expiry, 
			class, 
			registration_date,
			pucc_upto, 
			pucc_number,
            chassis_number, 
			engine_number, 
			fuel_type, 
			brand_name, 
			brand_model, 
			cubic_capacity, 
			gross_weight, 
			cylinders, 
			color, 
			norms,
			seating_capacity, 
			owner_count,
			fitness, 
			tax_upto, 
			permit_number, 
            permit_valid_upto, 
			permit_type, 
			national_permit_number, 
			national_permit_upto, 
			national_permit_issued_by, 
			total_challans, 
			pending_challans,
			rc_status
		)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	ChallanInsert = `
        INSERT INTO challans (
            challan_no, 
			license_plate, 
			date, 
			accused_name, 
			challan_status, 
			amount, 
			state, 
			area, 
			offence
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	DBInstance *sql.DB
)

func OpenDB() error {

	databasePath := os.Getenv("DB_PATH")

	fmt.Println("db path is valid ->", filepath.IsLocal(databasePath))

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1)
	DBInstance = db
	return err
}

func CreateDB() error {
	err := OpenDB()
	if err != nil {
		return err
	}

	tx, err := DBInstance.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to start transaction to create table")
	}

	stmt, err := tx.Prepare(`CREATE TABLE IF NOT EXISTS vehicles (
						license_plate TEXT PRIMARY KEY UNIQUE NOT NULL,
						owner_name TEXT,
						father_name TEXT,
						is_financed INTEGER,
						financer TEXT,
						present_address TEXT,
						permanent_address TEXT,
						insurance_company TEXT,
						insurance_policy TEXT,
						insurance_expiry TEXT,
						class TEXT,
						registration_date TEXT,
						vehicle_age TEXT,
						pucc_upto TEXT,
						pucc_number TEXT,
						chassis_number TEXT UNIQUE NOT NULL,
						engine_number TEXT UNIQUE NOT NULL,
						fuel_type TEXT,
						brand_name TEXT,
						brand_model TEXT,
						cubic_capacity TEXT,
						gross_weight TEXT,
						cylinders TEXT,
						color TEXT,
						norms TEXT,
						noc_details TEXT,
						seating_capacity TEXT,
						owner_count TEXT,
						fitness TEXT,
						tax_upto TEXT,
						tax_paid_upto TEXT,
						permit_number TEXT,
						permit_issue_date TEXT,
						permit_valid_from TEXT,
						permit_valid_upto TEXT,
						permit_type TEXT,
						national_permit_number TEXT,
						national_permit_upto TEXT,
						national_permit_issued_by TEXT,
						total_challans INTEGER, 
						pending_challans INTEGER,
						rc_status TEXT
					);
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = tx.Prepare(`
				CREATE TABLE IF NOT EXISTS challans (
					challan_no TEXT PRIMARY KEY UNIQUE,
					license_plate TEXT NOT NULL,
					date TEXT,
					accused_name TEXT,
					challan_status TEXT,
					amount TEXT,
					state TEXT,
					area TEXT,
					offence TEXT
				);
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt.Close()

	stmt, err = tx.Prepare(`
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY,
				name TEXT,
				email TEXT UNIQUE,
				password TEXT,
				access TEXT	
			);
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt.Close()

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit table creation")
	}

	return nil
}
