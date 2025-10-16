package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/manishbadgotra/vehicle-details/database"
	"github.com/manishbadgotra/vehicle-details/models"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))

func Login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Fprintf(w, "recovered from error: %v\n", recovered)
		}
	}()

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

	session, err := store.New(r, "session")
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
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}

	session.Save(r, w)

	json.NewEncoder(w).Encode(&user)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Fprintf(w, "recovered from error: %v\n", recovered)
		}
	}()

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
