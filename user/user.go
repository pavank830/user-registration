package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func checkEmailExists(email string) (bool, error) {
	_, err := getLoginInfoFromEmail(email)
	switch {
	case err == sql.ErrNoRows: // email  doesn't exists in db
		return false, nil
	case err != nil: // internal server error 500
		return false, err
	default: // email exists route to login page
		return true, nil
	}
}

func hashPasswordAndAddUser(userData User, id string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("server error, unable to create your account %v", err.Error())
	}
	userData.Password = string(hashedPassword)
	return addUserEntry(userData, id)
}

func validateuserLogin(email, password string) (string, int, error) {
	var id string
	userInfo, err := getLoginInfoFromEmail(email)
	if err != nil {
		return id, http.StatusInternalServerError, err
	}
	if userInfo == nil {
		return id, http.StatusOK, fmt.Errorf("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.HashedPassword), []byte(password))
	if err != nil {
		return id, http.StatusOK, err
	}
	id = userInfo.UserID
	return id, http.StatusOK, err
}

// CheckInBlackList -- chekc the jwt token in backlist table
func CheckInBlackList(token string) bool {
	fmt.Println("------>1")
	check, err := getFromBlackList(token)
	if err != nil {
		log.Println("error checking the token from blacklist table:", err)
		check = true
	}
	fmt.Println("------>1", check)
	return check
}

// getUserProfile -- get user profile from cache is exists, if not get from db and add to cache
func getUserProfile(id string) (*User, error) {
	var err error
	fmt.Println("------> id", id)
	user := getFromCache(id)
	if user == nil {
		user, err = getUserDataFromDB(id)
		if user != nil && err == nil {

			addToCache(id, *user)
		}
	}
	return user, err
}
