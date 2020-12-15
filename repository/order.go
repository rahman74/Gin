package repository

import (
	"database/sql"
	"log"
	"test/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ordersRepository struct {
	db *sqlx.DB
}

// OrdersRepository ...
type OrdersRepository interface {
	GetAllOrders(*sqlx.DB) ([]*model.Orders, error)
	CreateOrder(tx *sql.Tx, order model.Orders) (int, error)
	UpdateOrder(tx *sql.Tx, order model.Orders) (bool, error)
	DeleteOrder(tx *sql.Tx, orderID int)(bool, error)
	GetMaxPONumber() (string, error)
}

// NewOrdersRepository ...
func NewOrdersRepository(db *sqlx.DB) OrdersRepository {
	return &ordersRepository{
		db,
	}
}


func (ordersRepository *ordersRepository) GetMaxPONumber() (string, error){
	row := ordersRepository.db.QueryRow(`
		SELECT 
			order_number
		FROM public.orders
		ORDER BY order_id DESC LIMIT 1;
		`,
	)

	var orderNumber string

	sqlError := row.Scan(
		&orderNumber,
	)
	if sqlError != nil {
		log.Println("SQL error on GetMaxPONumber =>", sqlError)
	}

	return orderNumber, sqlError
}


// GetAllOrders ...
func (ordersRepository *ordersRepository) GetAllOrders(db *sqlx.DB) ([]*model.Orders, error) {
	var orders []*model.Orders
	rows, sqlError := db.Query(`
			SELECT order_id, customer_id, order_number, order_date, payment_method_id
			FROM public.orders;
		`,
	)

	if sqlError != nil {
		log.Print("SQL error on GetAllOrders => ", sqlError)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ID, customerID, paymentMethodID int
			var orderDate time.Time
			var orderNumber string

			sqlError := rows.Scan(
				&ID,
				&customerID,
				&orderNumber,
				&orderDate,
				&paymentMethodID,
			)

			if sqlError != nil {
				log.Print("SQL error on GetAllOrders => ", sqlError)
			} else {
				orders = append(
					orders,
					&model.Orders{
						CustomerID:          customerID,
						PaymentMethodID:        paymentMethodID,
						OrderID: ID,
						OrderNumber: orderNumber,
						OrderDate: orderDate,
					},
				)
			}
		}
	}

	if sqlError == sql.ErrNoRows {
		sqlError = nil
	}

	return orders, sqlError
}



func (ordersRepository *ordersRepository) CreateOrder(tx *sql.Tx, order model.Orders) (int, error) {
	var id int

	sqlError := tx.QueryRow(`
		INSERT INTO public.orders
		(customer_id, order_number, order_date, payment_method_id)
		VALUES($1, $2, now(), $3)	
		RETURNING order_id
		`,
		order.CustomerID,
		order.OrderNumber,
		order.PaymentMethodID,
	).Scan(&id)

	if sqlError != nil {
		log.Print("SQL error on CreateOrder => ", sqlError)
	}

	return id, sqlError
}


// UpdateOrder ...
func (ordersRepository *ordersRepository) UpdateOrder(tx *sql.Tx, order model.Orders) (bool, error) {
	var effectedRow int64

	result, sqlError := tx.Exec(`
	UPDATE public.orders
		SET payment_method_id=$1 where order_id = $2;

		`,
		order.PaymentMethodID,
		order.OrderID,
	)

	if sqlError != nil {
		log.Print("SQL error on UpdateOrder => ", sqlError)
	} else {
		effectedRow, sqlError = result.RowsAffected()
		if sqlError != nil {
			log.Print("SQL error on UpdateOrder => ", sqlError)
		}
	}

	return effectedRow > 0, sqlError
}


// DeleteOrder ...
func (ordersRepository *ordersRepository) DeleteOrder(tx *sql.Tx, orderID int) (bool, error) {
	var effectedRow int64

	result, sqlError := tx.Exec(`
	delete from public.orders where order_id = $1;
		`,
		orderID,
	)

	if sqlError != nil {
		log.Print("SQL error on DeleteOrder => ", sqlError)
	} else {
		effectedRow, sqlError = result.RowsAffected()
		if sqlError != nil {
			log.Print("SQL error on DeleteOrder => ", sqlError)
		}
	}

	return effectedRow > 0, sqlError
}
