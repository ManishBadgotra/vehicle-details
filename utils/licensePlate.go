package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
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

	for range time.Tick(time.Second) {
		// Open the CSV file
		log.Println("opening vehicles.csv file")

		file, err := os.Open("vehicles.csv")
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close() // Ensure the file is closed

		// Create a new CSV reader
		reader := csv.NewReader(file)
		log.Println("new csv reader created")

		fmt.Println("\nReading records line by line:")
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

			_, _, errResp := models.FetchVehicleDetails(payload)
			if errResp != nil {
				log.Printf("error in fetching `License Number: %v's` details with error: %v\n", licensePlate, errResp.Error)
			}
		}

		// log.Fatalf("whole csv data fetched")
	}

}
