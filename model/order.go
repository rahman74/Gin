package model

import "time"

//Orders Model
type Orders struct {
	OrderID int `json:"OrderID"`
	CustomerID int `json:"CustomerID"`
	OrderNumber string `json:"OrderNumber"`
	OrderDate time.Time `json:"OrderDate"`
	PaymentMethodID int `json:"PaymentMethodID"`
	OrderDetails []*OrderDetails `json:"OrderDetails"`
	PaymentMethod *PaymentMethods `json:"PaymentMethod"`
	Customer *Customer `json:"Customer"`
}

//OrderDetails ...
type OrderDetails struct {
	OrderDetailID int `json:"OrderDetailID"`
	ProductID int `json:"ProductID"`
	OrderID int `json:"OrderID"`
	QTY int `json:"QTY"`
	CreatedDate time.Time `json:"CreatedDate"`
	Product *Product `json:"Product"`
}

//Product ...
type Product struct {
	ProductID int `json:"ProductID"`
	ProductName string `json:"ProductName"`
	BasicPrice float64 `json:"BasicPrice"`
	CreatedDate time.Time `json:"CreatedDate"`
}

//PaymentMethods ...
type PaymentMethods struct {
	PaymentMethodID int `json:"PaymentMethodID"`
	MethodName string `json:"MethodName"`
	Code string `json:"Code"`
	CreatedDate time.Time `json:"CreatedDate"`
}