package database

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite" // this package is required for sqlite driver
)

var (
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

var (
	db   *sql.DB
	once sync.Once
	err  error
)

func OpenDB() (*sql.DB, error) {

	once.Do(func() {
		db, openError := sql.Open("sqlite", "your_database_file.db")
		if openError != nil {
			err = openError
			return
		}
		_, execError := db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			err = execError
			return
		}

		db.SetMaxOpenConns(1)
	}) // due to once.Do this will always returns a singleton instance of DB no matter how many times its called.

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
						license_plate TEXT PRIMARY KEY UNIQUE,
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
						vehicle_age INTEGER,
						pucc_upto TEXT,
						pucc_number TEXT,
						chassis_number TEXT UNIQUE,
						engine_number TEXT UNIQUE,
						fuel_type TEXT,
						brand_name TEXT,
						brand_model TEXT,
						cubic_capacity INTEGER,
						gross_weight INTEGER,
						cylinders INTEGER,
						color TEXT,
						norms TEXT,
						noc_details TEXT,
						seating_capacity INTEGER,
						owner_count INTEGER,
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
