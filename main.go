package main

import (
	"net/http"

	"github.com/ashwin-m/transactions/controllers/accounts"
	"github.com/ashwin-m/transactions/controllers/transactions"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func setupRoutes(r *gin.Engine) {

	// setup routes for accounts
	accountsHandler := accounts.NewHandler()
	accountsHandler.RouteGroup(r)

	// setup routes for transactions
	transactionsHandler := transactions.NewHandler()
	transactionsHandler.RouteGroup(r)
}

func main() {
	r := setupRouter()
	setupRoutes(r)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
