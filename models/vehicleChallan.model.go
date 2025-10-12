package models

type VehicleChallans struct {
	Code      int             `json:"code,omitempty"`
	Status    string          `json:"status,omitempty"`
	Message   string          `json:"message,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	Response  ChallanResponse `json:"response,omitempty"`
}

type ChallanResponse struct {
	RequestID string     `json:"request_id,omitempty"`
	VehicleID string     `json:"vehicleId,omitempty"`
	Total     int        `json:"total,omitempty"`
	Challans  []challans `json:"challans,omitempty"`
}

type challans struct {
	ChallanNo     string        `json:"challan_no,omitempty"`
	Date          string        `json:"date,omitempty"`
	AccusedName   string        `json:"accused_name,omitempty"`
	ChallanStatus string        `json:"challan_status,omitempty"`
	Amount        int           `json:"amount,omitempty"`
	State         string        `json:"state,omitempty"`
	Area          string        `json:"area,omitempty"`
	Offence       string        `json:"offence,omitempty"`
	OffenceList   []offenceList `json:"offence_list,omitempty"`
}

func NewChallan() *challans {
	return &challans{}
}

func NewOffenceList() []offenceList {
	var offenseList []offenceList
	return offenseList
}

type offenceList struct {
	OffenceName string `json:"offence_name,omitempty"`
}

func NewVehicleChallanResponse() *VehicleChallans {
	return &VehicleChallans{}
}
