package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pavank830/user-registration/utils"
)

// SignUp -- api to handle user signup
func SignUp(w http.ResponseWriter, r *http.Request) {

	var req SignUpReq
	resp := LoginResp{
		APIResp: APIResp{
			ResponseCode: utils.ResponseOK,
		},
	}
	var err error
	var emailCheck bool
	if err := utils.ParseRequest(w, r, &req); err != nil {
		return
	}

	defer func() {
		if err != nil {
			resp.ResponseCode = utils.ResponseFailed
			resp.ResponseDescription = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		}
		out, _ := json.Marshal(resp)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}()
	if req.UserData.Email == "" {
		err = errors.New("email is empty")
		return
	}
	if req.UserData.Password == "" {
		err = errors.New("Password is empty")
		return
	}
	emailCheck, err = checkEmailExists(req.UserData.Email)
	if err != nil {
		return
	}
	if emailCheck {
		http.ServeFile(w, r, "login.html")
		return
	}
	id := uuid.New()
	err = hashPasswordAndAddUser(req.UserData, id.String())
	if err != nil {
		return
	}

	jwtToken, err := utils.GenerateJWT(id.String())
	if err != nil {
		return
	}
	resp.ResponseDescription = "User created!"
	resp.JWTToken = jwtToken
	w.WriteHeader(http.StatusOK)
}

// Login - api to handle user login
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	resp := LoginResp{
		APIResp: APIResp{
			ResponseCode: utils.ResponseOK,
		},
	}
	var err error
	if err := utils.ParseRequest(w, r, &req); err != nil {
		return
	}

	defer func() {
		out, _ := json.Marshal(resp)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}()

	if req.Email == "" {
		err = errors.New("email is empty")
	}
	if req.Password == "" {
		err = errors.New("Password is empty")
	}
	if err != nil {
		resp.ResponseCode = utils.ResponseFailed
		resp.ResponseDescription = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, code, err := validateuserLogin(req.Email, req.Password)
	if err != nil {
		resp.ResponseCode = utils.ResponseFailed
		resp.ResponseDescription = err.Error()
		w.WriteHeader(code)
		return
	}
	jwtToken, err := utils.GenerateJWT(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.ResponseCode = utils.ResponseFailed
		resp.ResponseDescription = err.Error()
		return
	}
	resp.JWTToken = jwtToken
	w.WriteHeader(http.StatusOK)
}

// Logout - api to handle user logout
func Logout(w http.ResponseWriter, r *http.Request) {
	resp := APIResp{
		ResponseCode: utils.ResponseOK,
	}
	var err error
	defer func() {
		if err != nil {
			resp.ResponseCode = utils.ResponseFailed
			resp.ResponseDescription = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		out, _ := json.Marshal(resp)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}()

	err = addToBlackList(r.Header.Get(utils.HeaderJWT))
	if err != nil {
		return
	}
	resp.ResponseDescription = "user logged out"
}

// GetProfile - api to get user profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	resp := ProfileResp{
		APIResp: APIResp{
			ResponseCode: utils.ResponseOK,
		},
	}
	var err error
	var userInfo *User
	defer func() {
		if err != nil {
			resp.ResponseCode = utils.ResponseFailed
			resp.ResponseDescription = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		out, _ := json.Marshal(resp)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	}()

	userInfo, err = getUserProfile(r.Header.Get(utils.HeaderUserID))
	if err != nil {
		log.Println("error getting user profile data from db", err)
		return
	}
	if userInfo == nil {
		err = errors.New("user does not exist")
		log.Println(err)
		return
	}
	resp.User = *userInfo
}
