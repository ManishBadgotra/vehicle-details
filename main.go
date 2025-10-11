package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/manishbadgotra/vehicle-details/controllers"
	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/models"
)

func init() {
	if err := database.CreateDB(); err != nil {
		log.Fatalf("unable to create tables in database table: %v", err.Error())
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {

		secretCode := r.Header.Get("x-secret-code")
		if secretCode != "ManishIsAGenius" {
			errResponse := models.NewErrorResponse("secret key not setup correctly")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
	})

	mux.HandleFunc("GET /v1/vehicles", controllers.GetVehicle)
	mux.HandleFunc("POST /v1/vehicles", controllers.AddVehicle)
	mux.HandleFunc("PUT /v1/vehicles", controllers.UpdateVehicle)
	mux.HandleFunc("DELETE /v1/vehicles", controllers.DeleteVehicle)

	fmt.Println("Listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
