package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mini-e-wallet/config"
	"mini-e-wallet/domain/consumer"
	"mini-e-wallet/helpers"
	"mini-e-wallet/repository"
	consumer_sub "mini-e-wallet/service/consumer"
	"os"
)

func init() {
	viper.SetConfigFile(`../../config/config.json`)
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

	kafkaUsr := viper.GetString(`kafka.username`)
	kafkaPwd := viper.GetString(`kafka.password`)
	kafkaAddr := viper.GetString(`kafka.address`)
	kafkaRetry := viper.GetInt(`kafka.retry`)
	kafkaTimeout := viper.GetInt(`kafka.timeout`)
	kafkaConfig := config.GetKafkaConfig(kafkaUsr, kafkaPwd, kafkaTimeout, kafkaRetry)
	consumers, err := sarama.NewConsumer([]string{kafkaAddr}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Error create kakfa consumer got error %v", err)
	}
	defer func() {
		if err := consumers.Close(); err != nil {
			logrus.Fatal(err)
			return
		}
	}()

	kafkaConsumer := consumer.KafkaConsumer{
		Consumer: consumers,
	}

	walletRepo := repository.NewWalletRepository(db)
	walletService := consumer_sub.NewSubWalletService(kafkaConsumer, walletRepo)

	signals := make(chan os.Signal, 1)
	walletService.Process([]string{helpers.WalletTopic}, signals)
}
