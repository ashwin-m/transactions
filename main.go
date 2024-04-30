package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	accounts_controller "github.com/ashwin-m/transactions/controllers/accounts"
	"github.com/ashwin-m/transactions/controllers/transactions"
	accounts_dao "github.com/ashwin-m/transactions/daos/accounts"
	transactions_dao "github.com/ashwin-m/transactions/daos/transactions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func setupDB() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error is occurred  on .env file please check")
	}

	//we read our .env file
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, dbName)

	// set up postgres sql to open it.
	db, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return db
}

func setupRoutes(r *gin.Engine, dbPool *pgxpool.Pool, accountsDao accounts_dao.Dao, transactionsDao transactions_dao.Dao) {

	// setup routes for accounts
	accountsHandler := accounts_controller.NewHandler(accountsDao)
	accountsHandler.RouteGroup(r)

	// setup routes for transactions
	transactionsHandler := transactions.NewHandler(dbPool, accountsDao, transactionsDao)
	transactionsHandler.RouteGroup(r)
}

func main() {
	r := setupRouter()

	db := setupDB()
	defer db.Close()

	accountsDao := accounts_dao.NewDao(db)
	transactionsDao := transactions_dao.NewDao(db)

	setupRoutes(r, db, accountsDao, transactionsDao)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
