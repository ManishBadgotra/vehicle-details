package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/manishbadgotra/vehicle-details/database"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Access   string
}

func (u *UserRequest) LoginUser() (user UserRequest, err error) {
	conn, err := database.DBInstance.Conn(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error in creating connection")
		return UserRequest{}, err
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return UserRequest{}, err
	}

	result := tx.QueryRow(`
		SELECT name, email, password, access FROM users WHERE email == ?
	`,
		u.Email,
	)
	if result.Err() != nil {
		tx.Rollback()
		return UserRequest{}, result.Err()
	}

	if err = result.Scan(&user.Name, &user.Email, &user.Password, &user.Access); err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(os.Stderr, "no user found")
			tx.Rollback()
			return UserRequest{}, fmt.Errorf("no user found")
		} else {
			fmt.Fprintln(os.Stderr, err.Error()+" error in Scan function")
			tx.Rollback()
			return UserRequest{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return UserRequest{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return UserRequest{}, fmt.Errorf("password is incorrect")
	}

	return user, nil
}

func (u *UserRequest) CreateUser() error {
	conn, err := database.DBInstance.Conn(context.Background())
	if err != nil {
		return err
	}

	defer conn.Close()

	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to begin transaction")
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Access = "user"

	row := tx.QueryRow(`SELECT * FROM users WHERE email == ?`, u.Email)

	var newUser UserRequest
	if err := row.Scan(newUser.Email); err != nil {
		if err == sql.ErrNoRows {
			result := tx.QueryRow(`INSERT INTO users (name, email, password, access) VALUES (?, ?, ?, ?)`, u.Name, u.Email, hashedPassword, u.Access)
			if result.Err() != nil {
				fmt.Fprintf(os.Stderr, "error in inserting new user")
				tx.Rollback()
				return result.Err()
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("user already exist")
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Fprintf(os.Stderr, "error in commiting transaction")
		return err
	}

	return nil
}

// type UserSessionResponse struct {
// 	UserRequest
// 	Session string `json:"session"`
// }

// func UserSession(session string, user UserRequest) *UserSessionResponse {
// 	return &UserSessionResponse{
// 		UserRequest: user,
// 		Session:     session,
// 	}
// }
