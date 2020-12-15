package main

import (
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"test/handler"
	"test/repository"
	"test/usecase"
	"test/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const pubKeyPath = "keys/app.rsa.pub"

func main() {
	// Gin write log file
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	log.Println("service starting")

	// Get all config
	config := utils.GetConfig()

	// Setup log
	utils.InitializeLogSetting(&config.Log)

	log.Println("service initialization")

	// Get DB
	db := utils.ConnectDB(config)

	// Build Repositories
	customerRepository := repository.NewCustomerRepository(db)
	orderRepository := repository.NewOrdersRepository(db)
	orderDetailsRepository := repository.NewOrderDetailsRepository(db)
	productRepository := repository.NewProductRepository(db)
	paymentMethodsRepository := repository.NewPaymentMethodsRepository(db)

	// Build Usecases
	customerUsecase := usecase.NewCustomerUsecase(db, *config, config.Server, config.Jwt, customerRepository)
	orderUsecase := usecase.NewOrdersUsecase(db,orderDetailsRepository,orderRepository,customerRepository,productRepository,paymentMethodsRepository)

	// Build Handlers
	customerHandler := handler.NewCustomerHandler(config.Jwt, customerUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	// router := gin.Default()
	// Gin custom log format
	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	router.POST("/login", customerHandler.Login)
	router.POST("/customer", customerHandler.Customer)
	router.GET("/order", tokenAuthMiddleware() ,orderHandler.ReadOrder)
	router.POST("/add-order", tokenAuthMiddleware(), orderHandler.CreateOrder)
	router.DELETE("/delete-order/:id", tokenAuthMiddleware(), orderHandler.DeleteOrder)
	router.PUT("/update-order",tokenAuthMiddleware(), orderHandler.UpdateOrder)
	log.Println("service started on port ", config.Server.Port)
	router.Run(fmt.Sprintf(":%v", config.Server.Port))
}


func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var verifyKey *rsa.PublicKey
		var token *jwt.Token

		verifyBytes, err := ioutil.ReadFile(pubKeyPath)
		if err == nil {
			verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			jwtToken := c.Request.Header.Get("token-auth")

			if !strings.Contains(jwtToken, "Bearer") {
				err = fmt.Errorf("Invalid token")
				log.Printf(err.Error())
				c.JSON(http.StatusUnauthorized, err)
				c.Abort()
				return
			}

			jwtToken = strings.Replace(jwtToken,"Bearer ","",-1)
			// validate the token
			token, err = jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
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
		}

		if err == nil {
			c.Next()
		} else {
			log.Printf(err.Error())
			c.JSON(http.StatusUnauthorized, err)
			c.Abort()
			return
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
