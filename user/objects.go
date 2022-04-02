package user

// api paths
const (
	SignUpPath  = "/api/signup"
	LoginPath   = "/api/login"
	LogoutPath  = "/api/logout"
	ProfilePath = "/api/profile"
)

// User -- contains user basic info
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

//SignUpReq -- contains user basic info
type SignUpReq struct {
	UserData User `json:"user"`
}

//LoginReq -- contains email and password
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//LoginInfo -- contains email and password
type LoginInfo struct {
	UserID         string
	Email          string
	HashedPassword string
}

// APIResp - basic any api response struct
type APIResp struct {
	ResponseCode        int    `json:"response_code"`
	ResponseDescription string `json:"response_description"`
}

// ProfileResp -- resp to get user profile endpoint
type ProfileResp struct {
	User
	APIResp
}

// LoginResp -resp of login api
type LoginResp struct {
	JWTToken string `json:"token"`
	APIResp
}
