package usecase

import (
	"fmt"
	"test/model"
	"test/repository"
	"test/utils"

	"github.com/jmoiron/sqlx"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	accessToken = "level1"
	saltSize = 16

)

type customerUsecase struct {
	db                     *sqlx.DB
	config                 model.Config
	serverConfig           model.ServerConfig
	jwtConfig              model.JwtConfig
	customerRepository         repository.CustomerRepository
}

// CustomerUsecase ...
type CustomerUsecase interface {
	CreateCustomer(customer model.Customer) (int, error)
	GetCustomerByEmailOrPhoneNumber(string)(*model.Customer, error)
	AuthenticateUser(password string, user *model.Customer) (bool, error)}

// NewCustomerUsecase ...
func NewCustomerUsecase(
	db *sqlx.DB,
	config model.Config,
	serverConfig model.ServerConfig,
	jwtConfig model.JwtConfig,
	customerRepository repository.CustomerRepository,
) CustomerUsecase {
	return &customerUsecase{
		db,
		config,
		serverConfig,
		jwtConfig,
		customerRepository,
	}
}

func(customerUsecase *customerUsecase) CreateCustomer(customer model.Customer)(int, error) {
	var checkCustomer *model.Customer
	var err error
	var newID int
	tx, err := customerUsecase.db.Begin()
	checkCustomer, _ = customerUsecase.customerRepository.GetCustomerByEmailOrPhoneNumber(customer.Email, customer.PhoneNumber)
	if(checkCustomer.CustomerID == 0){
		salt := utils.GenerateRandomSalt(saltSize)
		passwordHashed := utils.HashPassword(customer.Password, salt)
		customer.Password = passwordHashed
		customer.Salt = salt
		newID, err = customerUsecase.customerRepository.CreateCustomer(tx, customer)
	}else{
		err = fmt.Errorf("customer already exists")
	}
	return newID, err
}

func(customerUsecase *customerUsecase) GetCustomerByEmailOrPhoneNumber(email string)(*model.Customer, error){
	return customerUsecase.customerRepository.GetCustomerByEmailOrPhoneNumber(email,email)
}

func (customerUsecase *customerUsecase) AuthenticateUser(password string, user *model.Customer) (bool, error) {
	var err error
	var isAuthenticated bool


	if (utils.DoPasswordsMatch(user.Password, password, user.Salt)){
		isAuthenticated = true
	} else {
		isAuthenticated = false
	}
	return isAuthenticated, err
}
