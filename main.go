package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

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

	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "x-request-code"},
	}

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {

		secretCode := r.Header.Get("x-request-code")
		if secretCode != "ManishIsAGenius" {
			errResponse := models.NewErrorResponse("secret key not setup correctly")
			json.NewEncoder(w).Encode(errResponse)
			return
		}

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
	})

	mux.HandleFunc("POST /v1/login", controllers.Login)
	mux.HandleFunc("POST /v1/signup", controllers.Signup)

	mux.HandleFunc("GET /v1/vehicles", controllers.GetVehicle)
	mux.HandleFunc("POST /v1/vehicles", controllers.AddVehicle)
	mux.HandleFunc("PUT /v1/vehicles", controllers.UpdateVehicle)
	mux.HandleFunc("DELETE /v1/vehicles", controllers.DeleteVehicle)

	fmt.Println("Listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", cors.New(corsOptions).Handler(mux)))
}
