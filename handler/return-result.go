package handler

import (
	"net/http"
	"test/model"

	"github.com/gin-gonic/gin"
)

// Create
func createResult(c *gin.Context, id int, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"ID":    id,
			"Error": "",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"ID":    id,
			"Error": err.Error(),
		})
	}
}

// Read
func readOrderResult(c *gin.Context, orders []*model.Orders, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"Orders": orders,
			"Error":        "",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Orders": orders,
			"Error":        err.Error(),
		})
	}
}
// Update Or Delete
func updateDeleteResult(c *gin.Context, isSuccess bool, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"IsSuccess": isSuccess,
			"Error":     "",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"IsSuccess": isSuccess,
			"Error":     err.Error(),
		})
	}
}