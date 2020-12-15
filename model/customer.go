package model

import "time"

//Customer Model
type Customer struct {
	CustomerID int `json:"CustomerID"`
	CustomerName string `json:"CustomerName"`
	Email string `json:"Email"`
	PhoneNumber string `json:"Phone_Number"`
	DOB time.Time `json:"DOB"`
	Sex string `json:Sex`
	Salt []byte `json:Salt`
	Password string `json:Password`
	CreatedDate time.Time `json:CreatedDate`
}

//CustomerJSONBody ...
type CustomerJSONBody struct {
	Action         string         `form:"Action" json:"Action" xml:"Action"`
	Customer           Customer           `form:"Customer" json:"Customer" xml:"Customer"`
}

// UserToken ...
type UserToken struct {
	ID                   int    `form:"ID" json:"ID" xml:"ID"`
	Email                string `form:"Email" json:"Email" xml:"Email"`
	Name                 string `form:"Name" json:"Name" xml:"Name"`
}