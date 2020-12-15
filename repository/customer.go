package repository

import (
	"database/sql"
	"log"
	"test/common"
	"test/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type customerRepository struct {
	db *sqlx.DB
}

// CustomerRepository ...
type CustomerRepository interface {
	CreateCustomer(tx *sql.Tx, customer model.Customer) (int, error)
	GetCustomerByEmailOrPhoneNumber(email, phone string) (*model.Customer, error)
	GetAllCustomer(db *sqlx.DB) ([]*model.Customer, error)
}

// NewCustomerRepository ...
func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRepository{
		db,
	}
}

func(customerRepository *customerRepository) CreateCustomer(tx *sql.Tx, customer model.Customer)(int, error) {
	var id int

	sqlError := tx.QueryRow(`
	INSERT INTO public.customers
	(customer_name, email, phone_number, dob, sex, salt, "password", created_date)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING customer_id
	
		`,
		customer.CustomerName,
		customer.Email,
		customer.PhoneNumber,
		customer.DOB,
		customer.Sex,
		customer.Salt,
		customer.Password,
		time.Now(),
	).Scan(&id)

	if sqlError != nil {
		log.Print("SQL error on CreateCustomer => ", sqlError)
	}
	common.TxChecker(tx, sqlError)
	return id, sqlError
}

func (customerRepository *customerRepository) GetCustomerByEmailOrPhoneNumber(emailValue, phone string) (*model.Customer, error){
	row := customerRepository.db.QueryRow(`
		SELECT 
			customer_id,
			customer_name,
			email,
			phone_number,
			dob,
			sex,
			salt,
			"password",
			created_date
		FROM public.customers
		WHERE email = $1 or phone_number = $2;	
		`,
		emailValue,
		phone,
	)

	var id int
	var name, email, phoneNumber, sex, password string
	var salt []byte
	var dob, createdDate time.Time

	sqlError := row.Scan(
		&id,
		&name,
		&email,
		&phoneNumber,
		&dob,
		&sex,
		&salt,
		&password,
		&createdDate,
	)

	if sqlError != nil {
		log.Println("SQL error on GetCustomerByEmailOrPhoneNumber =>", sqlError)
	}

	return &model.Customer{
			CustomerID:                   id,
			CustomerName:           name,
			Email:                email,
			PhoneNumber:              phoneNumber,
			DOB:                dob,
			Password:             password,
			Sex:              sex,
			Salt:                salt,
			CreatedDate:             createdDate,
		},
		sqlError
}


// GetAllCustomer ...
func (customerRepository *customerRepository) GetAllCustomer(db *sqlx.DB) ([]*model.Customer, error) {
	var customers []*model.Customer
	rows, sqlError := db.Query(`
			SELECT 
				customer_id,
				customer_name,
				email,
				phone_number,
				dob,
				sex,
				salt,
				"password",
				created_date
			FROM public.customers;
		`,
	)

	if sqlError != nil {
		log.Print("SQL error on GetAllCustomer => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name, email, phoneNumber, sex, password string
			var salt []byte
			var dob, createdDate time.Time

			sqlError := rows.Scan(
				&id,
				&name,
				&email,
				&phoneNumber,
				&dob,
				&sex,
				&salt,
				&password,
				&createdDate,
			)

			if sqlError != nil {
				log.Print("SQL error on GetAllCustomer => ", sqlError)
			} else {
				customers = append(
					customers,
					&model.Customer{
						CustomerID:                   id,
						CustomerName:           name,
						Email:                email,
						PhoneNumber:              phoneNumber,
						DOB:                dob,
						Password:             password,
						Sex:              sex,
						Salt:                salt,
						CreatedDate:             createdDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return customers, sqlError
}