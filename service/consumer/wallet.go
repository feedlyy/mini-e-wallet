package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/domain/consumer"
	"os"
	"time"
)

type subWalletService struct {
	kafka      consumer.KafkaConsumer
	walletRepo domain.WalletRepository
}

func NewSubWalletService(k consumer.KafkaConsumer, w domain.WalletRepository) consumer.SubWalletService {
	return &subWalletService{
		kafka:      k,
		walletRepo: w,
	}
}

func (w *subWalletService) Process(topics []string, signals chan os.Signal) {
	var (
		chanMessage   = make(chan *sarama.ConsumerMessage, 256)
		err           error
		partitionList []int32
		walletAcc     domain.Wallets
		ctx           = context.Background()
		tr            = domain.Transaction{}
	)
	for _, topic := range topics {
		partitionList, err = w.kafka.Consumer.Partitions(topic)
		if err != nil {
			logrus.Errorf("Sub-Service|Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(w.kafka.Consumer, topic, partition, chanMessage)
		}
	}
	logrus.Infof("Sub-Service|Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			err = json.Unmarshal(msg.Value, &tr)
			if err != nil {
				logrus.Errorf("Sub-Service|Failed marshal, err:%v", err)
				break ConsumerLoop
			}

			// wait for 5s
			time.Sleep(5 * time.Second)

			// get data wallet
			walletAcc, err = w.walletRepo.GetByOwnedID(ctx, tr.TransactionBy)
			if err != nil && err != sql.ErrNoRows {
				break ConsumerLoop
			}

			walletAcc.Balance += tr.Amount
			err = w.walletRepo.Update(ctx, walletAcc)
			if err != nil {
				break ConsumerLoop
			}
			logrus.Println("Wallet Updated!")
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		logrus.Errorf("Sub-Service|Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			logrus.Errorf("Sub-Service|Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}
}
