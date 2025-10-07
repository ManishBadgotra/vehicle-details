package models

type vehicleChallanResponse struct {
	Status bool            `json:"status,omitempty"`
	Result []ChallanResult `json:"data,omitempty"`
}
type ChallanResult struct {
	ChallanNo         string   `json:"challan_no,omitempty"`
	ChallanDate       string   `json:"challan_date,omitempty"`
	ChallanStatus     string   `json:"challan_status,omitempty"`
	Amount            string   `json:"amount,omitempty"`
	PaymentDate       any      `json:"payment_date,omitempty"`
	Offences          []string `json:"offences,omitempty"`
	ChallanPaymentURL string   `json:"challan_payment_url,omitempty"`
	ViolatorName      string   `json:"violator_name,omitempty"`
	DocNo             string   `json:"doc_no,omitempty"`
	ChallanReceipt    string   `json:"challan_receipt,omitempty"`
}

func NewVehicleChallanResponse() *vehicleChallanResponse {
	return &vehicleChallanResponse{}
}
