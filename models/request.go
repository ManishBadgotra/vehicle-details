package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

type request struct {
	RegNO       string `json:"vehicle_no"`
	Consent     string `json:"consent"`
	ConsentText string `json:"consent_text"`
	Config      Config
}

func NewVehicleRequest() *request {
	req := request{
		Consent:     "Y",
		ConsentText: "I hear by declare my consent agreement for fetching my information via AITAN Labs API",
	}
	return &req
}

type Config struct {
	URL              string
	APIKey           string
	VehicleEndpoint  string
	ChallansEndpoint string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("no .env file found")
	}

	return &Config{
		URL:              os.Getenv("URL"),
		APIKey:           os.Getenv("API_KEY"),
		VehicleEndpoint:  os.Getenv("VEHICLE_ENDPOINT"),
		ChallansEndpoint: os.Getenv("CHALLAN_ENDPOINT"),
	}
}

func (request *request) GetVehicleDetails(w http.ResponseWriter, r *http.Request) {
	requestURL := path.Join(request.Config.URL, request.Config.VehicleEndpoint)

	fmt.Println(requestURL)
	secretCode := r.Header.Get("x-secret-code")
	if secretCode != "ManishIsAGenius" {
		errResponse := NewErrorResponse("secret key not setup correctly")
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusNoContent)

		errResp := NewErrorResponse("something went related to your request")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	request.RegNO = strings.Trim(request.RegNO, " ")
	if request.RegNO == "" {
		w.WriteHeader(http.StatusExpectationFailed)

		errResp := NewErrorResponse("no vehicle number provided")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	payload, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := NewErrorResponse("unable to create response for vehicle number")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payload))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := NewErrorResponse("unable to make request to the server")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req.Header.Add("x-rapidapi-key", request.Config.APIKey)
	req.Header.Add("x-rapidapi-host", "rto-vehicle-information-india.p.rapidapi.com")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch res.StatusCode {
	case http.StatusOK:

		vehicleStruct := NewVehicleDetailsResponse()

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))

		// Save body to file
		err = os.WriteFile("vehicle_details.json", body, 0644)
		if err != nil {
			log.Println("Failed to save response:", err)
		}

		if err = json.Unmarshal(body, &vehicleStruct); err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)

		fmt.Println("Vehicle Details Request successfull")
		json.NewEncoder(w).Encode(vehicleStruct)

	case http.StatusTooManyRequests:

		w.WriteHeader(http.StatusTooManyRequests)
		errResp := NewErrorResponse("request quota exceeded.")
		json.NewEncoder(w).Encode(errResp)

	default:

		w.WriteHeader(res.StatusCode)
		errResp := NewErrorResponse(fmt.Sprintf("error occured on third party request: %d", res.StatusCode))
		json.NewEncoder(w).Encode(errResp)

	}
}

func (request *request) GetVehicleChallans(w http.ResponseWriter, r *http.Request) {

	requestURL := path.Join(request.Config.URL, request.Config.ChallansEndpoint)

	secretCode := r.Header.Get("x-secret-code")
	if secretCode != "ManishIsAGenius" {
		errResponse := NewErrorResponse("secret key not setup correctly")
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		w.WriteHeader(http.StatusNoContent)

		errResp := NewErrorResponse("something went wrong related to your request")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	request.RegNO = strings.Trim(request.RegNO, " ")
	if request.RegNO == "" {
		w.WriteHeader(http.StatusExpectationFailed)

		errResp := NewErrorResponse("no vehicle number provided")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	payload, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := NewErrorResponse("unable to create response for vehicle number")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payload))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errResp := NewErrorResponse("unable to make request to the server")
		json.NewEncoder(w).Encode(errResp)
		return
	}

	req.Header.Add("x-rapidapi-key", request.Config.APIKey)
	req.Header.Add("x-rapidapi-host", "rto-vehicle-information-india.p.rapidapi.com")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		vehicleStruct := NewVehicleChallanResponse()

		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		fmt.Println(string(body))

		// Save body to file
		err = os.WriteFile("vehicle_challans.json", body, 0644)
		if err != nil {
			log.Println("Failed to save response:", err)
		}

		if err = json.Unmarshal(body, &vehicleStruct); err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println("Vehicle Challans Request successfull")
		json.NewEncoder(w).Encode(vehicleStruct)

	case http.StatusTooManyRequests:

		w.WriteHeader(http.StatusTooManyRequests)
		errResp := NewErrorResponse("request quota exceeded.")
		json.NewEncoder(w).Encode(errResp)

	default:

		w.WriteHeader(res.StatusCode)
		errResp := NewErrorResponse(fmt.Sprintf("error occured on third party request: %d", res.StatusCode))
		json.NewEncoder(w).Encode(errResp)

	}
}
