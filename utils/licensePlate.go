package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/manishbadgotra/vehicle-details/models"
)

func VerifyVehicleNumber(licensePlate string) bool {
	var (
		isValid      bool
		vehicleRegex string = `^(?:[A-Z]{2}[0-9]{1,2}[A-Z]{1,3}[0-9]{4}|[0-9]{1,2}BH[0-9]{4}[A-Z]{1,2})$`
	)

	regex, err := regexp.Compile(vehicleRegex)
	if err != nil {
		log.Printf("error in vehicle ID Verify regex, error: %v", err)
		return false
	}

	isValid = regex.Match([]byte(licensePlate))

	return isValid
}

type vehicleStruct struct {
	VehicleId string `json:"vehicleId"`
}

func GetVehiclesFromList() {

	for t := range time.Tick(time.Second) {

		day := os.Getenv("WEEKDAY")

		fmt.Println(t)
		fmt.Println(t.Date())
		fmt.Println(t.Weekday())
		fmt.Println(day)

		if strings.EqualFold(t.Weekday().String(), day) {
			// Open the CSV file
			// log.Println("opening vehicles.csv file")

			file, err := os.Open("vehicles.csv")
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer file.Close() // Ensure the file is closed

			// Create a new CSV reader
			reader := csv.NewReader(file)
			// log.Println("new csv reader created")

			fmt.Println("\nReading vehicle line by line:")
			for {
				record, err := reader.Read()
				if err == io.EOF {
					break // End of file
				}
				if err != nil {
					log.Fatalf("Error reading record: %v", err)
				}

				licensePlate := record[0]

				reqVehicle := vehicleStruct{}
				reqVehicle.VehicleId = licensePlate

				payload, err := json.Marshal(&reqVehicle)
				if err != nil {
					log.Printf("License Number: %v is unable to marshal", licensePlate)
					return
				}

				_, statusCode, errResp := models.FetchVehicleDetails(payload)
				if errResp != nil {

					forDay := 24 * time.Hour

					if statusCode == 400 {
						log.Printf("bad request for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 401 {
						log.Printf("Unauthorized/Expired for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 402 {
						log.Printf("Insufficient Funds stopping server for few hours and sending alert on mail \n")
					}

					if statusCode == 403 {
						log.Printf("Unauthenticated Request while requesting data for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 404 {
						log.Printf("Not Found for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 405 {
						log.Printf("Method Not Allowed for `License Number: %v\n", reqVehicle.VehicleId)
						os.Exit(1)
					}

					if statusCode == 415 {
						log.Printf("Unsupported Media Type for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 422 {
						log.Printf("Request failed due to invalid details for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 429 {
						log.Printf("Too many requests for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 500 {
						log.Printf("Internal Server Error while fetching data for `License Number: %v\n", reqVehicle.VehicleId)
					}

					if statusCode == 503 {
						log.Printf("Backend Down/Maintenance stopping server for few hours for `License Number: %v\n", reqVehicle.VehicleId)
					}

					log.Println("cron job on sleep for a day")
					time.Sleep(forDay)

				}

				if statusCode == http.StatusOK {
					log.Printf("Reqeust successfull for `License Number: %v`\n", reqVehicle.VehicleId)
				}

			}
		}

		// log.Fatalf("whole csv data fetched")
	}

}
