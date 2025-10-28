package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/models"
)

var store = sessions.NewCookieStore([]byte{101, 99, 79, 97, 87, 216, 73, 16, 114, 201, 31, 208, 155, 56, 110, 225, 196, 141, 136, 104, 140, 149, 146, 67, 168, 234, 17, 216, 226, 10, 56, 29, 174, 37, 66, 236, 141, 157, 96, 237, 199, 188, 168, 13, 254, 54, 241, 79, 216, 94, 27, 157, 70, 187, 39, 11, 120, 209, 167, 169, 117, 172, 43, 50})

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.UserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		fmt.Fprintf(w, "recovered from error: %v\n", err)
	}

	user, err := user.LoginUser()
	if err != nil {

		fmt.Fprintf(os.Stderr, "time: %v - error in function Login User: %v\n", time.Now(), err.Error())

		w.Header().Add("Content-Type", "application/json")
		errResp := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(errResp)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	session, err := store.New(r, "user-session")
	if err != nil {
		fmt.Fprintf(os.Stderr, "time: %v - error in creating user: %v\n", time.Now(), err.Error())

		w.Header().Add("Content-Type", "application/json")
		errResp := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 1,
		HttpOnly: true,
		Secure:   false,
	}

	session.Save(r, w)

	json.NewEncoder(w).Encode(&user)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if err := database.DBInstance.Ping(); err != nil {
		panic("error unable to ping db")
	}

	var user models.UserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		fmt.Fprintf(w, "recovered from error: %v\n", err)
	}

	if err := user.CreateUser(); err != nil {
		fmt.Println(err.Error())

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)

		errResp := models.NewErrorResponse(err.Error())
		json.NewEncoder(w).Encode(errResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}
