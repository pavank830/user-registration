package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// basic api resp codes
const (
	ResponseOK     = 0
	ResponseFailed = 1
)

// http custom headers
var (
	HeaderJWT    = "jwt"
	HeaderUserID = "user_id"
)

// JWTSecretKey -
var JWTSecretKey = []byte("jwtsecretkey")

// jwt token errors
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// ParseRequest -- func to read and parse the http request/input
func ParseRequest(w http.ResponseWriter, r *http.Request, input interface{}) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body %+v\n", err)
		return err
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("error closing body %s\n", err.Error())
		return err
	}
	if err := json.Unmarshal(body, input); err != nil {
		log.Printf("Unmarshalling Error. %+v ", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		return err
	}
	return nil
}

// GenerateJWT - to generate jwt token
func GenerateJWT(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
