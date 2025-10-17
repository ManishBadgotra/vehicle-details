package utils

import "regexp"

func VerifyVehicleNumber(licensePlate string) (bool, error) {
	var (
		isValid      bool
		vehicleRegex string = `^(?:[A-Z]{2}[0-9]{1,2}[A-Z]{1,3}[0-9]{4}|[0-9]{1,2}BH[0-9]{4}[A-Z]{1,2})$`
	)

	regex, err := regexp.Compile(vehicleRegex)

	isValid = regex.Match([]byte(licensePlate))

	return isValid, err
}
