package handler

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"test/model"
	"test/usecase"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	privKeyPath = "keys/app.rsa" //openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" //openssl rsa -in app.rsa > app.rsa.pub
	accessToken = "level1"
)

type customerHandler struct {
	jwtConfig             model.JwtConfig
	customerUsecase       usecase.CustomerUsecase
}

// CustomerHandler ....
type CustomerHandler interface {
	Customer(c *gin.Context)
	Login(c *gin.Context)
}

// NewCustomerHandler ...
func NewCustomerHandler(
	jwtConfig model.JwtConfig,
	customerUsecase  usecase.CustomerUsecase,
) CustomerHandler {
	return &customerHandler{
		jwtConfig,
		customerUsecase,
	}
}

func (customerHandler *customerHandler) Customer(c *gin.Context) {
	var customerJSONBody model.CustomerJSONBody
	var err error

	if err = c.ShouldBindJSON(&customerJSONBody); err == nil {
		switch strings.ToUpper(customerJSONBody.Action) {
		case createCustomer:
			customerHandler.CreateCustomer(c, customerJSONBody.Customer)
		}
	}
	fmt.Println(err)
}

func (customerHandler *customerHandler) CreateCustomer(c *gin.Context, customer model.Customer) {
	id, err := customerHandler.customerUsecase.CreateCustomer(customer)
	createResult(c, id, err)
}

func (customerHandler *customerHandler) Login(c *gin.Context) {
	var loginJSON model.Login
	var isAuthorized bool
	var err error
	var token string
	var loginResult model.LoginResult
	var isRefreshToken bool

	if err = c.ShouldBindJSON(&loginJSON); err == nil {
		switch strings.ToUpper(loginJSON.Action) {
		case refreshToken:
			isRefreshToken = true
			loginResult, _ = customerHandler.jwtRsaDecryptToken(loginJSON, isRefreshToken)
		default:
			if loginJSON.Email != "" && loginJSON.Password != "" {
				loginJSON.JwtToken = c.Request.Header.Get("Login-Token")
				if loginJSON.JwtToken == "" {
					loginResult = customerHandler.checkLogin(loginJSON)
				} else {
					loginResult, _ = customerHandler.jwtRsaDecryptToken(loginJSON, isRefreshToken)
				}
			} else {
				loginResult = setLoginResult(http.StatusBadRequest, isAuthorized, token, "Username and password can't be empty.")
			}
		}
	} else {
		loginResult = setLoginResult(http.StatusBadRequest, isAuthorized, token, err.Error())
	}

	c.JSON(
		loginResult.HTTPStatus,
		gin.H{
			"IsAuthorized": loginResult.IsAuthorized,
			"Token":        loginResult.Token,
			"ErrorMessage": loginResult.ErrorMessage,
		},
	)
}



func (customerHandler *customerHandler) jwtRsaDecryptToken(loginJSON model.Login, isRefreshToken bool) (model.LoginResult, error) {
	var loginResult model.LoginResult
	var verifyKey *rsa.PublicKey

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err == nil {
		verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

		// validate the token
		token, err := jwt.Parse(loginJSON.JwtToken, func(token *jwt.Token) (interface{}, error) {
			// since we only use the one private key to sign the tokens,
			// we also only use its public counter part to verify
			return verifyKey, nil
		})

		// branch out into the possible error from signing
		switch err.(type) {
		case nil: // no error
			if !token.Valid { // but may still be invalid
				err = fmt.Errorf("Invalid token")
			}

		case *jwt.ValidationError: // something was wrong during the validation
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				err = fmt.Errorf("Token Expired")

			default:
				err = fmt.Errorf("Token ValidationError error: %v", vErr.Errors)
			}

		default: // something else went wrong
			err = fmt.Errorf("Token parse error: %v", err)
		}

		if err == nil {
			loginResult = setLoginResult(http.StatusOK, true, loginJSON.JwtToken, "")
		} else if err.Error() == "Token Expired" || isRefreshToken {
			err = nil
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				var signKey *rsa.PrivateKey
				var newTokenStr string
				signBytes, err := ioutil.ReadFile(privKeyPath)
				if err == nil {
					signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
					if err == nil {
						// create a signer for rsa 256
						newToken := jwt.New(jwt.GetSigningMethod("RS256"))

						// set the expire time
						claims["exp"] = time.Now().Add(time.Minute * time.Duration(customerHandler.jwtConfig.ExpirationDurationInMinute)).Unix()

						newToken.Claims = claims
						newTokenStr, err = token.SignedString(signKey)
						if err == nil {
							loginJSON.JwtToken = newTokenStr
							fmt.Println(newTokenStr)
							loginResult = setLoginResult(http.StatusOK, true, loginJSON.JwtToken, "")
						} else {
							err = fmt.Errorf("Token Signing error: %v", err)
						}
					}
				}
			}

			if err != nil {
				loginResult = setLoginResult(http.StatusBadRequest, false, loginJSON.JwtToken, err.Error())
			}
		} else {
			loginResult = setLoginResult(http.StatusBadRequest, false, loginJSON.JwtToken, err.Error())
		}
	}

	return loginResult, err
}

func setLoginResult(httpStatus int, isAuthorized bool, token, errorMessage string) model.LoginResult {
	return model.LoginResult{
		HTTPStatus:   httpStatus,
		IsAuthorized: isAuthorized,
		Token:        token,
		ErrorMessage: errorMessage,
	}
}

func (customerHandler *customerHandler) checkLogin(loginJSON model.Login) model.LoginResult {
	var isAuthorized bool
	var err error
	var customer *model.Customer
	var token string
	var loginResult model.LoginResult

	customer, err = customerHandler.customerUsecase.GetCustomerByEmailOrPhoneNumber(loginJSON.Email)
	if err == nil {

		isAuthorized, err = customerHandler.customerUsecase.AuthenticateUser(loginJSON.Password, customer)
		

		if err == nil && isAuthorized {
			token, err = customerHandler.jwtRsaGenerateToken(loginJSON, customer)
			if err == nil {
				loginResult = setLoginResult(http.StatusOK, isAuthorized, token, "")
			}
		} else {
			log.Println(fmt.Sprintf("User %v login error : %v", loginJSON.Email, err))
			err = fmt.Errorf("Unauthorized")
			loginResult = setLoginResult(http.StatusUnauthorized, isAuthorized, token, err.Error())
		}
	}

	return loginResult
}


///jwtRsaGenerateToken ..
func (customerHandler *customerHandler) jwtRsaGenerateToken(
	loginJSON model.Login,
	user *model.Customer,
) (string, error) {
	var signKey *rsa.PrivateKey
	var tokenString string

	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err == nil {
		signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err == nil {
			// create a signer for rsa 256
			token := jwt.New(jwt.GetSigningMethod("RS256"))

			// set our claims
			claims := make(jwt.MapClaims)
			claims["AccessToken"] = accessToken
			claims["CustomUserInfo"] = struct {
				Name string
				User          *model.UserToken
			}{
				loginJSON.Email,
				&model.UserToken{
					ID:                   user.CustomerID,
					Email:                user.Email,
					Name:                 user.CustomerName,
				},
			}

			// set the expire time
			claims["exp"] = time.Now().Add(time.Minute * time.Duration(customerHandler.jwtConfig.ExpirationDurationInMinute)).Unix()

			token.Claims = claims
			tokenString, err = token.SignedString(signKey)
			if err != nil {
				err = fmt.Errorf("Token Signing error: %v", err)
			}
		}
	}

	return tokenString, err
}