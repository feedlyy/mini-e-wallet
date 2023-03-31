package main

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	accHandler "mini-e-wallet/handler"
	md "mini-e-wallet/middleware"
	"mini-e-wallet/repository"
	"mini-e-wallet/service"
	"net/http"
	"time"
)

//go:embed migrations
var migrations embed.FS

func init() {
	viper.SetConfigFile(`../config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Setup Logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.name`)
	dbPort := viper.GetString(`database.port`)
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbName, dbPort)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Pong from db")

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	timeoutCtx := viper.GetInt(`context.timeout`)
	accountRepo := repository.NewAccountRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	accountService := service.NewAccountService(accountRepo, tokenRepo)
	walletService := service.NewWalletService(walletRepo, tokenRepo)

	accountHandler := accHandler.NewAccountHandler(accountService, time.Duration(timeoutCtx)*time.Second)
	walletHandler := accHandler.NewWalletHandler(walletService, time.Duration(timeoutCtx)*time.Second)
	middleware := md.NewMiddleware(tokenRepo)

	serverPort := viper.GetString(`server.address`)
	handler := httprouter.New()
	handler.POST("/api/v1/init", accountHandler.RegistUser)
	handler.POST("/api/v1/wallet", middleware.AuthMiddleware(walletHandler.EnableWallet))
	handler.PATCH("/api/v1/wallet", middleware.AuthMiddleware(walletHandler.DisableWallet))
	handler.GET("/api/v1/wallet", middleware.AuthMiddleware(walletHandler.ViewBalance))

	logrus.Infof("Server run on localhost%v", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, handler))
}
