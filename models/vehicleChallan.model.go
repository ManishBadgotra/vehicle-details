package models

type ChallanResponse struct {
	Code     int      `json:"code"`
	Status   string   `json:"status"`
	Message  string   `json:"message"`
	Response Response `json:"response"`
}

type Response struct {
	RequestID   string     `json:"request_id"`
	VehicleID   string     `json:"vehicleId"`
	Total       int        `json:"total"`
	ChallanList []Challans `json:"challans"`
}

type Challans struct {
	ChallanNo     string `json:"challan_no"`
	Date          string `json:"date"`
	AccusedName   string `json:"accused_name"`
	ChallanStatus string `json:"challan_status"`
	Amount        string `json:"amount"`
	State         string `json:"state"`
	Area          string `json:"area"`
	Offence       string `json:"offence"`
}

// type Offence struct {
// 	OffenceName string `json:"offence_name"`
// }

// type OffenceList struct {
// 	OffenceName string `json:"offence_name,omitempty"`
// }

func NewVehicleChallanResponse() *ChallanResponse {
	return &ChallanResponse{}
}
