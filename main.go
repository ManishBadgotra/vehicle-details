package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manishbadgotra/vehicle-details/models"
)

func main() {

	mux := http.NewServeMux()

	config := models.NewConfig()

	request := models.NewVehicleRequest()
	request.Config = *config

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {

		secretCode := r.Header.Get("x-secret-code")
		if secretCode != "ManishIsAGenius" {
			errResponse := models.NewErrorResponse("secret key not setup correctly")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		if request.Config.URL == "" || request.Config.APIKey == "" || request.Config.VehicleEndpoint == "" || request.Config.ChallansEndpoint == "" {
			errResp := models.NewErrorResponse("env not setup correctly")
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
	})

	mux.HandleFunc("POST /details", request.GetVehicleDetails)
	mux.HandleFunc("POST /challans", request.GetVehicleChallans)

	fmt.Println("Listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
