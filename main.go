package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "natsClient: ", 0)

	logger.Println("starting the application...")

	viper.SetConfigFile("config.yml") // Path to your configuration file
	if err := viper.ReadInConfig(); err != nil {
		logger.Println("error reading config file:", err)
		return
	}

	// Read the nats URL value from the configuration
	url := viper.GetString("nats.url")
	if url == "" {
		url = "nats://localhost:4222"
		logger.Println("No value found for nats ur, so using ", url)
	}

	// messages per second from config
	messagesPerSecondStr := viper.GetString("test.messagesPerSecond")
	if messagesPerSecondStr == "" {
		messagesPerSecondStr = "1"
		logger.Println("no value found for test:messagesPerSecond, so using ", messagesPerSecondStr)
	}

	messagesPerSecond, err := strconv.Atoi(messagesPerSecondStr)
	if err != nil {
		logger.Fatal(err)
	}

	// Connect to NATS server
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Fatalf("nats error:%s , \n URL: %s", url, err)
	}
	defer nc.Close()

	err = sendAndReceiveMessages(logger, nc, messagesPerSecond)
	if err != nil {
		logger.Fatal(err)
	}

	keepAlive(logger, err)

	logger.Println("exiting the application...")

}

func keepAlive(logger *log.Logger, err error) {
	duration := viper.GetString("test.duration")
	if duration == "" {
		duration = "1"
		logger.Println("no value found for test:duration, so using ", duration)
	}
	testDuration, err := strconv.Atoi(duration)
	ticker := time.NewTicker(time.Second * time.Duration(testDuration))

	// Infinite loop to keep the application running
	select {
	case <-ticker.C:
		os.Exit(0)
	}
}

func sendAndReceiveMessages(logger *log.Logger, nc *nats.Conn, messagesPerSecond int) error {
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// send messages
	go func() {
		ticker := time.NewTicker(time.Second * 1)

		publishCount := 1
		subject := "FX_DETERMINATION"

		// Infinite loop to continuously send messages for a given message count
		for range ticker.C {
			for i := 0; i < messagesPerSecond; i++ {
				publishMessage(logger, js, subject, publishCount)
				publishCount++
			}
		}

	}()

	// receive messages
	go func() {
		receiveCount := 1
		subject := "FX_CONVERTED"
		// Receive message from NATS subject
		_, err := js.Subscribe(subject, func(msg *nats.Msg) {
			logger.Printf("%s received msg: %d, %v\n", currentTime(), receiveCount, string(msg.Data))

			err := msg.Ack()
			if err != nil {
				log.Fatal(err)
				return
			}
			receiveCount++
		}, nats.AckExplicit(),
			nats.MaxAckPending(100))

		if err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func publishMessage(logger *log.Logger, js nats.JetStreamContext, subject string, publishCount int) {
	message := ConversionReq{
		TenantID:       12345,
		BankID:         123456,
		BaseCurrency:   "GBP",
		TargetCurrency: "USD",
		Tier:           "2",
		Amount:         randomFloat(),
		RequestedOn:    currentTime(),
	}
	msgByteArray, err := json.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}

	// Publish a message
	_, err = js.PublishAsync(subject, msgByteArray)
	if err != nil {
		log.Fatal(err)
	}
	logger.Printf("%s published msg: %d, %s \n", currentTime(), publishCount, string(msgByteArray))
}
