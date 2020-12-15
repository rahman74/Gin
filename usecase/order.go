package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"test/common"
	"test/model"
	"test/repository"
	"test/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

// orderUsecase ...
type ordersUsecase struct {
	db                              *sqlx.DB
	orderDetailsRepository           repository.OrderDetailsRepository
	ordersRepository       repository.OrdersRepository
	customerRepository repository.CustomerRepository
	productRepository repository.ProductRepository
	paymentMethodsRepository repository.PaymentMethodsRepository
}

// OrdersUsecase ...
type OrdersUsecase interface {
	GetAllOrders()([]*model.Orders, error)
	CreateOrders(order model.Orders) (int, error)
	DeleteOrder(orderID int) (bool, error)
	UpdateOrders(order model.Orders) (bool, error)
}

// NewOrdersUsecase ...
func NewOrdersUsecase(
	db                              *sqlx.DB,
	orderDetailsRepository           repository.OrderDetailsRepository,
	ordersRepository       repository.OrdersRepository,
	customerRepository repository.CustomerRepository,
	productRepository repository.ProductRepository,
	paymentMethodsRepository repository.PaymentMethodsRepository,
) OrdersUsecase {
	return &ordersUsecase{
		db,
		orderDetailsRepository,
		ordersRepository,
		customerRepository,
		productRepository,
		paymentMethodsRepository,
	}
}


func (ordersUsecase *ordersUsecase) GetAllOrders() ([]*model.Orders, error) {
	var err error
	var customers []*model.Customer
	 var products []*model.Product
	var paymentMethods []*model.PaymentMethods
	var orderDetailsTemp []*model.OrderDetails
	var orders []*model.Orders

	customersMap := make(map[int]*model.Customer)
	 productsMap := make(map[int]*model.Product)
	paymentMethodsMap := make(map[int]*model.PaymentMethods)
	// read order
	orders, err = ordersUsecase.ordersRepository.GetAllOrders(ordersUsecase.db)
	if err == nil {
		 // read product
		products, _ = ordersUsecase.productRepository.GetAllProduct(ordersUsecase.db)
		if err == nil {
			for _, product := range products {
				productsMap[product.ProductID] = product
			}
		}

		// read customer
		customers, _ = ordersUsecase.customerRepository.GetAllCustomer(ordersUsecase.db)
		if err == nil {
			for _, customer := range customers {
				customersMap[customer.CustomerID] = customer
			}
		}

		// read paymentMethods
		paymentMethods, _ = ordersUsecase.paymentMethodsRepository.GetAllPaymentMethods(ordersUsecase.db)
		if err == nil {
			for _, paymentMethod := range paymentMethods {
				paymentMethodsMap[paymentMethod.PaymentMethodID] = paymentMethod
			}
		}
		

		for _, order := range orders {
			// set customers
			if len(customers) > 0 {
				order.Customer = customersMap[order.CustomerID]
			}

			// set paymentMethods
			if len(paymentMethods) > 0 {
				order.PaymentMethod = paymentMethodsMap[order.PaymentMethodID]
			}
			orderDetailsTemp, _ = ordersUsecase.orderDetailsRepository.GetOrderDetailsByOrderID(order.OrderID, ordersUsecase.db)
			for _, orderDetail := range orderDetailsTemp {
				orderDetail.Product = productsMap[orderDetail.ProductID]
			}
			order.OrderDetails = orderDetailsTemp

		}
	}

	return orders, err
}

func (ordersUsecase *ordersUsecase) UpdateOrders(order model.Orders) (bool, error) {
	var isEdited bool

	tx, err := ordersUsecase.db.Begin()
	if err == nil {
		// update orders
		isEdited, err = ordersUsecase.ordersRepository.UpdateOrder(tx, order)
		if err == nil && order.OrderDetails != nil {
			for _, orderDetail := range order.OrderDetails {
				if orderDetail.OrderDetailID > 0 {
					// update order details
					_, err = ordersUsecase.orderDetailsRepository.UpdateOrderDetails(tx,orderDetail)
				} else {
					// create order details
					_, err = ordersUsecase.orderDetailsRepository.CreateOrderDetails(tx, order.OrderID,orderDetail)
				}
			}
		}

		common.TxChecker(tx, err)
	}

	return isEdited, err
}

func (ordersUsecase *ordersUsecase) DeleteOrder(orderID int) (bool, error) {
	var isDeleted bool

	tx, err := ordersUsecase.db.Begin()
	if err == nil {
		isDeleted, err = ordersUsecase.ordersRepository.DeleteOrder(tx, orderID)
		if err == nil {
			isDeleted, err = ordersUsecase.orderDetailsRepository.DeleteOrderDetails(tx, orderID)
		}
		common.TxChecker(tx, err)
	}

	return isDeleted, err
}


func (ordersUsecase *ordersUsecase) CreateOrders(order model.Orders) (int, error) {
	var newOrderID int

	tx, err := ordersUsecase.db.Begin()
	if err == nil {
		// create order
		order.OrderNumber = ordersUsecase.GenerateOrderNumber()
		newOrderID, err = ordersUsecase.ordersRepository.CreateOrder(tx, order)
		if err == nil && order.OrderDetails != nil {
			// create order details
			for _, orderDetail := range order.OrderDetails {
				_, err = ordersUsecase.orderDetailsRepository.CreateOrderDetails(tx, newOrderID, orderDetail)
			}
		}

		common.TxChecker(tx, err)
	}

	return newOrderID, err
}



func (ordersUsecase *ordersUsecase) GenerateOrderNumber()(string){
	var numberConv = 0
	var DigitNumber = 1
	lastNumber,_ := ordersUsecase.ordersRepository.GetMaxPONumber()
	year, month, _ := time.Now().Date()
	if(lastNumber != ""){
		splited := strings.Split(lastNumber, "/")
		splitedString := string(splited[0])
		splitedMonth := string(splited[1])
		splitedYear := string(splited[2])
		splitedYearInt, _ := strconv.Atoi(splitedYear)

		if(splitedMonth == utils.Roman(int(month)) && splitedYearInt == year){
			splitedNumber := strings.Split(splitedString,"-")
			numberConv, _ = strconv.Atoi(splitedNumber[1])
		}
		DigitNumber = numberConv + 1
		return fmt.Sprintf("PO-%d/%s/%d", DigitNumber, utils.Roman(int(month)), year)
	}
	return fmt.Sprintf("PO-%d/%s/%d",DigitNumber, utils.Roman(int(month)), year)
}