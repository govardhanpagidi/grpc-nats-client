package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/govardhanpagidi/nats-client/fxconversion"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {

	viper.SetConfigFile("config.yml") // Path to your configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Read the URL value from the configuration
	url := viper.GetString("nats.url")
	if url == "" {
		url = "nats://localhost:4222"
		fmt.Println("No value found for nats ur, so using ", url)
	}

	// messages per second
	messagesPerSecond := viper.GetString("test.messagesPerSecond")
	if messagesPerSecond == "" {
		messagesPerSecond = "1"
		fmt.Println("No value found for test:messagesPerSecond, so using ", messagesPerSecond)
	}

	msgsPerSecond, err := strconv.Atoi(messagesPerSecond)
	if err != nil {
		log.Fatal(err)
	}
	// Connect to NATS server
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	//err = sendMessageToNats(nc, subject)
	err = sendMessageToNatsJetStream(nc, msgsPerSecond)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func sendMessageToNatsJetStream(nc *nats.Conn, messagesPerSecond int) error {
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 1)

		publishCount := 1
		subject := "FX_DETERMINATION"

		// Infinite loop to continuously send messages
		for range ticker.C {
			for i := 0; i < messagesPerSecond; i++ {
				publishMessage(js, subject, publishCount)
				publishCount++
			}
		}

	}()

	go func() {
		receiveCount := 1
		subject := "FX_CONVERTED"
		// Receive message from NATS
		_, err := js.Subscribe(subject, func(msg *nats.Msg) {
			fmt.Printf("Received: msg: %d with %v\n", receiveCount, string(msg.Data))

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

func publishMessage(js nats.JetStreamContext, subject string, publishCount int) {
	message := ConversionReq{
		TenantID:       12345,
		BankID:         123456,
		BaseCurrency:   "GBP",
		TargetCurrency: "USD",
		Tier:           "2",
		Amount:         randomFloat(),
		RequestedOn:    time.Now().Format("2006-01-02T15:04:05.000"),
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
	fmt.Printf("Published: msg: %d with %s \n", publishCount, string(msgByteArray))
}

func sendMessageToNats(nc *nats.Conn, subject string) error {
	yourMessage := &pb.ConversionRequest{
		TenantID:       12347,
		BankID:         123456,
		BaseCurrency:   "GBP",
		TargetCurrency: "USD",
		Tier:           "2",
		Amount:         1000,
	}

	// Marshal your Protocol Buffer message
	protoBytes, err := proto.Marshal(yourMessage)
	if err != nil {
		log.Fatal(err)
	}

	// Publish gRPC message to NATS subject
	err = nc.Publish(subject, protoBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Published gRPC message to NATS subject:", subject)

	// This code is optional and only needed when you want to listen to a response
	_, err = nc.Subscribe("FX_CONVERTED", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))

		var response pb.ConversionResponse
		err := proto.Unmarshal(msg.Data, &response)
		if err != nil {
			log.Printf("Error unmarshalling gRPC message: %v\n", err)
			return
		}
		fmt.Printf("Received a message: %v\n", response)
	})

	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func randomFloat() float64 {
	rand.Seed(time.Now().UnixNano())

	// Define the range for the random float (e.g., between 1.0 and 10.0)
	minVal := 1.0
	maxVal := 100000.0

	// Generate a random float within the specified range
	return roundToTwoDecimalPlaces(minVal + rand.Float64()*(maxVal-minVal))
}

func roundToTwoDecimalPlaces(num float64) float64 {
	return float64(int(num*100)) / 100
}
