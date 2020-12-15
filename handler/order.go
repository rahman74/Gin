package handler

import (
	"strconv"
	"test/model"
	"test/usecase"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderUsecase usecase.OrdersUsecase
}

// OrderHandler ...
type OrderHandler interface {
	CreateOrder(c *gin.Context)
	ReadOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

// NewOrderHandler ...
func NewOrderHandler(orderUsecase usecase.OrdersUsecase) OrderHandler {
	return &orderHandler{
		orderUsecase,
	}
}

func (orderHandler *orderHandler) CreateOrder(c *gin.Context) {
	var orderJSON model.Orders
	err := c.ShouldBindJSON(&orderJSON)
	id, err := orderHandler.orderUsecase.CreateOrders(orderJSON)
	createResult(c, id, err)
}

func (orderHandler *orderHandler) ReadOrder(c *gin.Context) {
	apps, err := orderHandler.orderUsecase.GetAllOrders()
	readOrderResult(c, apps, err)
}

func (orderHandler *orderHandler) UpdateOrder(c *gin.Context) {
	var orderJSON model.Orders
	err := c.ShouldBindJSON(&orderJSON)
	isSuccess, err := orderHandler.orderUsecase.UpdateOrders(orderJSON)
	updateDeleteResult(c, isSuccess, err)
}

func (orderHandler *orderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	isSuccess, err := orderHandler.orderUsecase.DeleteOrder(idInt)
	updateDeleteResult(c, isSuccess, err)
}
