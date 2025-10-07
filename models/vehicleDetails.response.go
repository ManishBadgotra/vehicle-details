package models

type vehicleDetailsResponse struct {
	Status bool   `json:"status,omitempty"`
	Result Result `json:"data"`
}
type Result struct {
	RegistrationAuthority string `json:"registration_authority,omitempty"`
	RegistrationNo        string `json:"registration_no,omitempty"`
	RegistrationDate      string `json:"registration_date,omitempty"`
	ChassisNo             string `json:"chassis_no,omitempty"`
	EngineNo              string `json:"engine_no,omitempty"`
	OwnerName             string `json:"owner_name,omitempty"`
	VehicleClass          string `json:"vehicle_class,omitempty"`
	FuelType              string `json:"fuel_type,omitempty"`
	MakerModel            string `json:"maker_model,omitempty"`
	FitnessUpto           string `json:"fitness_upto,omitempty"`
	InsuranceUpto         string `json:"insurance_upto,omitempty"`
	FuelNorms             string `json:"fuel_norms,omitempty"`
	VehicleInfo           any    `json:"vehicle_info,omitempty"`
	VehicleType           string `json:"vehicle_type,omitempty"`
	InsuranceCompany      string `json:"insurance_company,omitempty"`
	FinancierName         any    `json:"financier_name,omitempty"`
	PucUpto               string `json:"puc_upto,omitempty"`
	RoadTaxPaidUpto       string `json:"road_tax_paid_upto,omitempty"`
	VehicleColor          string `json:"vehicle_color,omitempty"`
	SeatCapacity          string `json:"seat_capacity,omitempty"`
	UnloadWeight          string `json:"unload_weight,omitempty"`
	BodyTypeDesc          string `json:"body_type_desc,omitempty"`
	ManufactureMonthYear  string `json:"manufacture_month_year,omitempty"`
	RcStatus              string `json:"rc_status,omitempty"`
	Ownership             string `json:"ownership,omitempty"`
	OwnershipDesc         string `json:"ownership_desc,omitempty"`
}

func NewVehicleDetailsResponse() *vehicleDetailsResponse {
	return &vehicleDetailsResponse{}
}
