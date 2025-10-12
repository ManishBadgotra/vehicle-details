package database

import (
	"database/sql"

	_ "modernc.org/sqlite" // this package is required for sqlite driver
)

var (
	FindInVehicleTable = `
		SELECT id, license_plate, owner_name, father_name, is_financed, financer, present_address, permanent_address,
		insurance_company, insurance_policy, insurance_expiry, class, registration_date, vehicle_age, pucc_upto, pucc_number,
		chassis_number, engine_number, fuel_type, brand_name, brand_model, cubic_capacity, gross_weight, cylinders, color, norms,
		noc_details, seating_capacity, owner_count, tax_upto, tax_paid_upto, permit_number, permit_issue_date, permit_valid_from,
		permit_valid_upto, permit_type, national_permit_number, national_permit_upto, national_permit_issued_by, rc_status FROM vehicles 
		WHERE license_plate = ?
   		OR chassis_number = ?
   		OR engine_number = ? 
		LIMIT 1
	`
	FindInChallansTable = `
		SELECT id, challan_no, date, accused_name, challan_status, amount, state, area, offence, offence_list FROM challans WHERE license_plate = ?
	`
	VehicleInsert = `
        INSERT INTO vehicles
            (request_id, license_plate, owner_name, father_name, is_financed, financer, present_address, permanent_address,
            insurance_company, insurance_policy, insurance_expiry, class, registration_date, vehicle_age, pucc_upto, pucc_number,
            chassis_number, engine_number, fuel_type, brand_name, brand_model, cubic_capacity, gross_weight, cylinders, color, norms,
            noc_details, seating_capacity, owner_count, tax_upto, tax_paid_upto, permit_number, permit_issue_date, permit_valid_from,
            permit_valid_upto, permit_type, national_permit_number, national_permit_upto, national_permit_issued_by, rc_status)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	ChallanInsert = `
        INSERT INTO challans (
            challan_no, license_plate, date, accused_name, challan_status, amount, state, area, offence, offence_list
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	return db, err
}

func CreateDB() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS vehicles (
						request_id TEXT,
						license_plate TEXT PRIMARY KEY UNIQUE NOT NULL,
						owner_name TEXT,
						father_name TEXT,
						is_financed TEXT,
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

	stmt, err = db.Prepare(`
				CREATE TABLE IF NOT EXISTS challans (
					challan_no TEXT PRIMARY KEY,
					license_plate TEXT NOT NULL,
					date TEXT,
					accused_name TEXT,
					challan_status TEXT,
					amount INTEGER,
					state TEXT,
					area TEXT,
					offence TEXT,
					offence_list TEXT,
					FOREIGN KEY (license_plate) REFERENCES vehicles(license_plate)
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

	return nil
}
