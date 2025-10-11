package models

type requestBody struct {
	LicensePlate  string `json:"vehicleId"`
	ChassisNumber string `json:"chassis"`
	EngineNumber  string `json:"engine_no"`
}

func NewRequestBody(vehicleId, chassis, engine string) *requestBody {
	return &requestBody{
		LicensePlate:  vehicleId,
		ChassisNumber: chassis,
		EngineNumber:  engine,
	}
}
