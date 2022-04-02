package userregistration

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pavank830/user-registration/user"
	"github.com/pavank830/user-registration/utils"
)

func init() {
	tr := &http.Transport{
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   3,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       configTLS(""),
	}

	httpClient = &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}

//Start -- func that starts http server
func Start(port string, dsn string) error {
	user.DSN = dsn
	return startHTTPSvr(port)
}

func configTLS(serverName string) *tls.Config {
	TLSConfig := &tls.Config{
		ServerName:             serverName,
		Certificates:           []tls.Certificate{{}},
		Rand:                   rand.Reader,
		SessionTicketsDisabled: false,
		MinVersion:             tls.VersionTLS12,
		InsecureSkipVerify:     true,
	}
	return TLSConfig
}

func getRouter() *mux.Router {
	// gorilla/mux for routing
	router := mux.NewRouter()
	router.Use(starterMiddleware)
	router.HandleFunc(user.SignUpPath, user.SignUp).Methods(http.MethodPost)
	router.HandleFunc(user.LoginPath, user.Login).Methods(http.MethodPost)
	router.HandleFunc(user.LogoutPath, user.Logout).Methods(http.MethodPost)
	router.HandleFunc(user.ProfilePath, user.GetProfile).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(defaultHandler)
	return router
}

func starterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		if _, ok := NoAuth[req.RequestURI]; !ok {
			log.Println("jwt token verification")
			authHeader := strings.Split(req.Header.Get("Authorization"), "Bearer ")
			id, code, errmsg := validateJWTToken(authHeader)
			if errmsg != "" || code != 0 {
				resp.WriteHeader(code)
				resp.Write([]byte(errmsg))
				return
			}
			req.Header.Set(utils.HeaderUserID, id)
			if len(authHeader) > 0 {
				req.Header.Set(utils.HeaderJWT, authHeader[1])
			}
			fmt.Println("------> id", id)
			fmt.Println("------> id", authHeader[1])
			// check if token blacklisted
			if user.CheckInBlackList(authHeader[1]) {
				resp.WriteHeader(http.StatusUnauthorized)
				resp.Write([]byte(utils.ErrInvalidToken.Error()))
				return
			}
		}
		next.ServeHTTP(resp, req)
		log.Printf("Duration: %s %s %s\n", req.Method, req.RequestURI, time.Since(start).String())
	})
}

func startHTTPSvr(port string) error {
	server := &http.Server{
		Addr:         listenAddr + ":" + port,
		Handler:      getRouter(),
		ReadTimeout:  time.Duration(HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(HTTPWriteTimeout) * time.Second,
	}
	log.Printf("## HTTP Server listening on %v\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start http server %+v\n", err)
		return err
	}
	return nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("defaultHandler Duration: %s %s ", r.Method, r.RequestURI)
}

func validateJWTToken(authHeader []string) (string, int, string) {
	var code int
	var errmsg string
	var id string
	if len(authHeader) != 2 {
		log.Println("invalid token.")
		code = http.StatusUnauthorized
		errmsg = utils.ErrInvalidToken.Error()
		return id, code, errmsg
	}
	jwtToken := authHeader[1]
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utils.JWTSecretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorExpired {
				log.Println("token expired")
				code = http.StatusUnauthorized
				errmsg = utils.ErrExpiredToken.Error()
				return id, code, errmsg
			}
		}
		log.Println("token invalid,unauthorized")
		code = http.StatusUnauthorized
		errmsg = utils.ErrInvalidToken.Error()
		return id, code, errmsg
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("token invalid,claims empty")
		code = http.StatusUnauthorized
		errmsg = utils.ErrInvalidToken.Error()
		return id, code, errmsg
	}

	switch idt := claim["id"].(type) {
	case string:
		id = idt
		fmt.Println("------> id 10", idt, id)
	default:
		log.Println("token invalid,format")
		code = http.StatusUnauthorized
		errmsg = utils.ErrInvalidToken.Error()
		return id, code, errmsg
	}
	return id, code, errmsg
}
