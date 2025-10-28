package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/manishbadgotra/vehicle-details/models"
)

var (
	apiKey = ""
)

type vehicleStruct struct {
	VehicleId string `json:"vehicleId"`
}

func GetVehicleDetails(w http.ResponseWriter, r *http.Request) {

	licensePlate := r.URL.Query().Get("license")

	slog.String("Path Params --> ", licensePlate)

	existingVehicle := models.VehicleRequest{}
	vehicle, err := existingVehicle.GetFromDB(licensePlate)
	if err != nil {
		errResp := models.NewErrorResponse("data not found")
		json.NewEncoder(w).Encode(errResp)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(vehicle.Response)

}

func GetAllVehicleDetails(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("user-session")

	if err != nil {
		if err == http.ErrNoCookie {
			// Cookie not found
			fmt.Fprintf(w, "Cookie 'myCookie' not found.")
			return
		}
		// Other error occurred
		http.Error(w, "Error retrieving cookie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	t := cookie.Expires

	fmt.Fprintln(w, "all vehicles here --- ", "Value of 'myCookie': ", cookie.Value, " t is ---> ", t.Before(time.Now().UTC()))
}
