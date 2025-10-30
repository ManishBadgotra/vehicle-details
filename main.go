package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/manishbadgotra/vehicle-details/controllers"
	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/utils"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	godotenv.Load()

	if err := database.CreateDB(); err != nil {
		log.Fatalf("unable to create tables in database table: %v", err.Error())
		return
	}
}

func main() {

	go utils.GetVehiclesFromList()

	mux := http.NewServeMux()

	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}

	mux.HandleFunc("POST /v1/login", controllers.Login)
	mux.HandleFunc("POST /v1/signup", controllers.Signup)

	mux.HandleFunc("GET /v1/vehicles", controllers.GetVehicle)
	mux.HandleFunc("POST /v1/vehicles", controllers.AddVehicle)
	// mux.HandleFunc("PUT /v1/vehicles", controllers.UpdateVehicle)
	// mux.HandleFunc("DELETE /v1/vehicles", controllers.DeleteVehicle)

	// mux.HandleFunc("GET /v1/challans", controllers.GetVehicleChallans)

	log.Println("running PORT on :5898")
	log.Fatal(http.ListenAndServe(":5898", cors.New(corsOptions).Handler(mux)))
}
