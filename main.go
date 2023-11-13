package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/govardhanpagidi/nats-client/fxconversion"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
)

func main() {

	viper.SetConfigFile("config.yml") // Path to your configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// Read the URL value from the configuration
	url := viper.GetString("nats:url")
	if url == "" {
		url = "nats://localhost:4222"
		fmt.Println("No value found for nats ur, so using ", url)
	}

	// Connect to NATS server
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subject for gRPC messages
	subject := "FX_DETERMINATION"

	//err = sendMessageToNats(nc, subject)
	err = sendMessageToNatsJetStream(nc, subject)
	if err != nil {
		log.Fatal(err)
	}

}

func sendMessageToNatsJetStream(nc *nats.Conn, subject string) error {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Marshal your Protocol Buffer message
	message := ConversionReq{
		TenantID:       12347,
		BankID:         123456,
		BaseCurrency:   "GBP",
		TargetCurrency: "USD",
		Tier:           "2",
		Amount:         1000,
	}
	msgByteArray, err := json.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}

	// Publish gRPC message to NATS subject
	_, err = js.Publish(subject, msgByteArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Published message to nats jetstream:", string(msgByteArray))
	fmt.Println(" subject:", subject)

	//// This code is optional and only needed when you want to listen to a response
	//_, err = nc.Subscribe("FX_CONVERTED", func(msg *nats.Msg) {
	//	fmt.Printf("Received a message: %s\n", string(msg.Data))
	//
	//	var response pb.ConversionResponse
	//	err := proto.Unmarshal(msg.Data, &response)
	//	if err != nil {
	//		log.Printf("Error unmarshalling gRPC message: %v\n", err)
	//		return
	//	}
	//
	//	fmt.Printf("Received a message: %v\n", response)
	//
	//})
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	return nil
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
