package user_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	userregistration "github.com/pavank830/user-registration"
	"github.com/pavank830/user-registration/user"
	"github.com/pavank830/user-registration/utils"
)

const (
	host = "http://127.0.0.1:60010"
)

var signUpTest = []struct {
	user       user.User
	statusCode int
}{{
	user: user.User{
		FirstName: "pavan",
		LastName:  "kumar",
		Email:     "kumar830@gmail.com",
		Password:  "pavan",
	},
	statusCode: http.StatusOK,
}, {
	user: user.User{
		FirstName: "pavan",
		LastName:  "kumar",
		Email:     "pavank830@gmail.com",
		Password:  "",
	},
	statusCode: http.StatusInternalServerError,
},
}

var loginTest = []struct {
	user       user.LoginReq
	statusCode int
	code       int //api resp code
}{{
	user: user.LoginReq{
		Email:    "kumar830@gmail.com",
		Password: "pavan",
	},
	statusCode: http.StatusOK,
	code:       0,
}, {
	user: user.LoginReq{
		Email:    "kumar830@gmail.com",
		Password: "pavan2",
	},
	statusCode: http.StatusOK,
	code:       1,
},
}

var getProfileTest = []struct {
	JWTToken   string
	statusCode int
}{
	{
		JWTToken:   "",
		statusCode: http.StatusUnauthorized,
	}, {
		JWTToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDkwMTI3MzUsImlhdCI6MTY0ODkyNjMzNSwiaWQiOiIxZjc5YjFjMi1mZTdkLTRjNjktYTJkNS04YmM3NTFiNjlhZWYifQ.6hyaMQN41767GzM4XxK1XAvjf0jofbHYcILOOH75ZcU",
		statusCode: http.StatusUnauthorized,
	}, {
		JWTToken:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDkwMTI4MzMsImlhdCI6MTY0ODkyNjQzMywiaWQiOiIxZjc5YjFjMi1mZTdkLTRjNjktYTJkNS04YmM3NTFiNjlhZWYifQ.FvyDOvV-s6VKwe2jW-f8cPs5ATNvXRFOkhHIzKeGUVU",
		statusCode: http.StatusOK,
	},
}

func initdb() {
	user.DSN = "kumar:kumar@tcp(127.0.0.1:3307)/registration"
}
func TestSignUp(t *testing.T) {
	initdb()
	endpoint := host + user.SignUpPath
	for _, testData := range signUpTest {
		data, _ := json.Marshal(user.SignUpReq{
			UserData: testData.user,
		})
		req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
		w := httptest.NewRecorder()
		user.SignUp(w, req)
		res := w.Result()
		defer res.Body.Close()

		t.Logf("Expected sign up status code %d, got %d", testData.statusCode, res.StatusCode)
		if testData.statusCode != res.StatusCode {
			t.Errorf("Check Code. status code is not correct.")
			t.Fail()
		}
	}
}

func TestLogin(t *testing.T) {
	initdb()
	endpoint := host + user.LoginPath
	for _, testData := range loginTest {
		data, _ := json.Marshal(testData.user)
		req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
		w := httptest.NewRecorder()
		user.Login(w, req)
		res := w.Result()
		defer res.Body.Close()
		resp := user.LoginResp{}
		data, _ = ioutil.ReadAll(res.Body)
		err := json.Unmarshal(data, &resp)
		if err != nil {
			t.Fail()
		}
		t.Logf("Expected login api resp code %d, got %d", testData.code, resp.APIResp.ResponseCode)
		if testData.code != resp.APIResp.ResponseCode {
			t.Errorf("Check Code. login api resp code is not correct.")
			t.Fail()
		}
	}
}

func TestGetProfile(t *testing.T) {
	initdb()
	endpoint := host + user.ProfilePath
	for _, testData := range getProfileTest {
		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		w := httptest.NewRecorder()
		authHeader := []string{testData.JWTToken, testData.JWTToken}
		id, code, errmsg := userregistration.ValidateJWTToken(authHeader)
		if errmsg != "" || code != 0 {
			w.WriteHeader(code)
			w.Write([]byte(errmsg))
		} else {
			req.Header.Set(utils.HeaderUserID, id)
			if len(authHeader) > 0 {
				req.Header.Set(utils.HeaderJWT, authHeader[1])
			}
			user.GetProfile(w, req)
		}
		res := w.Result()
		defer res.Body.Close()

		t.Logf("Expected get profile status code %d, got %d", testData.statusCode, res.StatusCode)
		if testData.statusCode != res.StatusCode {
			t.Errorf("Check Code.  get profile is not correct.")
			t.Fail()
		}
	}
}

func TestLogout(t *testing.T) {
	initdb()
	endpoint := host + user.LogoutPath
	for _, testData := range getProfileTest {
		req := httptest.NewRequest(http.MethodPost, endpoint, nil)
		w := httptest.NewRecorder()
		authHeader := []string{testData.JWTToken, testData.JWTToken}
		id, code, errmsg := userregistration.ValidateJWTToken(authHeader)
		if errmsg != "" || code != 0 {
			w.WriteHeader(code)
			w.Write([]byte(errmsg))
		} else {
			req.Header.Set(utils.HeaderUserID, id)
			if len(authHeader) > 0 {
				req.Header.Set(utils.HeaderJWT, authHeader[1])
			}
			user.Logout(w, req)
		}
		res := w.Result()
		defer res.Body.Close()
		t.Logf("Expected logout status code %d, got %d", testData.statusCode, res.StatusCode)
		if testData.statusCode != res.StatusCode {
			t.Errorf("Check Code.logout is not correct.")
			t.Fail()
		}
	}
}
