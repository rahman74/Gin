package model

// Login ...
type Login struct {
	Email    string `form:"Email" json:"Email" xml:"Email"`
	Password string `form:"Password" json:"Password" xml:"Password"`
	JwtToken string `form:"JwtToken" json:"JwtToken" xml:"JwtToken"`
	Action   string `form:"Action" json:"Action" xml:"Action"`
}

// LoginResult ...
type LoginResult struct {
	IsAuthorized bool
	Token        string
	ErrorMessage string
	HTTPStatus   int
}

// TokenDetails ...
type TokenDetails struct {
	RefreshToken string
	AccessUuID   string
	RefreshUuID  string
	AtExpires    int64
	RtExpires    int64
}