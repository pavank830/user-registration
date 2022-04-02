package userregistration

import (
	"net/http"

	"github.com/pavank830/user-registration/user"
)

var (
	httpClient *http.Client
)

// http related constants
const (
	listenAddr       = ""
	HTTPReadTimeout  = 60
	HTTPWriteTimeout = 60
)

// NoAuth -endpoints for which jwt verifcation not needed
var NoAuth = map[string]struct{}{
	user.LoginPath:  {},
	user.SignUpPath: {},
}
