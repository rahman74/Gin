package repository

import (
	"database/sql"
	"log"
	"test/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type orderDetailsRepository struct {
	db *sqlx.DB
}

// OrderDetailsRepository ...
type OrderDetailsRepository interface {
	GetAllOrderDetails(*sqlx.DB) ([]*model.OrderDetails, error)
	GetOrderDetailsByOrderID(int,*sqlx.DB) ([]*model.OrderDetails, error)
	UpdateOrderDetails(tx *sql.Tx, orderDetails *model.OrderDetails) (bool, error)
	CreateOrderDetails(tx *sql.Tx, orderID int, orderDetails *model.OrderDetails) (int, error)
	DeleteOrderDetails(tx *sql.Tx, orderID int) (bool, error)
}

// NewOrderDetailsRepository ...
func NewOrderDetailsRepository(db *sqlx.DB) OrderDetailsRepository {
	return &orderDetailsRepository{
		db,
	}
}

// GetAllOrderDetails ...
func (orderDetailsRepository *orderDetailsRepository) GetAllOrderDetails(db *sqlx.DB) ([]*model.OrderDetails, error) {
	var orderDetails []*model.OrderDetails
	rows, sqlError := db.Query(`
			SELECT order_detail_id, order_id, product_id, qty, created_date
			FROM public.order_details;
		`,
	)

	if sqlError != nil {
		log.Print("SQL error on GetAllOrderDetails => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ID, orderID, productID, qty int
			var createdDate time.Time

			sqlError := rows.Scan(
				&ID,
				&orderID,
				&productID,
				&qty,
				&createdDate,
			)

			if sqlError != nil {
				log.Print("SQL error on GetAllOrderDetails => ", sqlError)
			} else {
				orderDetails = append(
					orderDetails,
					&model.OrderDetails{
						ProductID:          productID,
						OrderDetailID:        ID,
						OrderID: orderID,
						QTY: qty,
						CreatedDate: createdDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return orderDetails, sqlError
}

// GetOrderDetailsByOrderID ...
func (orderDetailsRepository *orderDetailsRepository) GetOrderDetailsByOrderID(ID int,db *sqlx.DB) ([]*model.OrderDetails, error) {
	var orderDetails []*model.OrderDetails
	rows, sqlError := db.Query(`
			SELECT order_detail_id, order_id, product_id, qty, created_date
			FROM public.order_details where order_id = $1;
		`,
		ID,
	)

	if sqlError != nil {
		log.Print("SQL error on GetOrderDetailsByOrderID => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ID, orderID, productID, qty int
			var createdDate time.Time

			sqlError := rows.Scan(
				&ID,
				&orderID,
				&productID,
				&qty,
				&createdDate,
			)

			if sqlError != nil {
				log.Print("SQL error on GetOrderDetailsByOrderID => ", sqlError)
			} else {
				orderDetails = append(
					orderDetails,
					&model.OrderDetails{
						ProductID:          productID,
						OrderDetailID:        ID,
						OrderID: orderID,
						QTY: qty,
						CreatedDate: createdDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return orderDetails, sqlError
}


// CreateOrderDetails ...
func (orderDetailsRepository *orderDetailsRepository) CreateOrderDetails(tx *sql.Tx, orderID int ,orderDetails *model.OrderDetails) (int, error) {
	var id int

	sqlError := tx.QueryRow(`
			INSERT INTO public.order_details
			(order_id, product_id, qty, created_date)
			VALUES($1, $2, $3, now())
			returning order_detail_id
		`,
		orderID,
		orderDetails.ProductID,
		orderDetails.QTY,
	).Scan(&id)

	if sqlError != nil {
		log.Print("SQL error on CreateOrderDetails => ", sqlError)
	}

	return id, sqlError
}


// UpdateOrderDetails ...
func (orderDetailsRepository *orderDetailsRepository) UpdateOrderDetails(tx *sql.Tx, orderDetails *model.OrderDetails) (bool, error) {
	var effectedRow int64

	result, sqlError := tx.Exec(`
	UPDATE public.order_details
		SET qty=$1 where order_detail_id = $2;

		`,
		orderDetails.QTY,
		orderDetails.OrderDetailID,
	)

	if sqlError != nil {
		log.Print("SQL error on UpdateOrderDetails => ", sqlError)
	} else {
		effectedRow, sqlError = result.RowsAffected()
		if sqlError != nil {
			log.Print("SQL error on UpdateOrderDetails => ", sqlError)
		}
	}

	return effectedRow > 0, sqlError
}

// DeleteOrderDetails ...
func (orderDetailsRepository *orderDetailsRepository) DeleteOrderDetails(tx *sql.Tx, orderID int) (bool, error) {
	var effectedRow int64

	result, sqlError := tx.Exec(`
	delete from public.order_details where order_id = $1;
		`,
		orderID,
	)

	if sqlError != nil {
		log.Print("SQL error on DeleteOrderDetails => ", sqlError)
	} else {
		effectedRow, sqlError = result.RowsAffected()
		if sqlError != nil {
			log.Print("SQL error on DeleteOrderDetails => ", sqlError)
		}
	}

	return effectedRow > 0, sqlError
}
