package main

import (
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
		fmt.Println("No value found for URL: using ", url)
	}

	// Connect to NATS server
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subject for gRPC messages
	subject := "FX_DETERMINATION"

	yourMessage := &pb.ConversionRequest{
		TenantID:       123,
		BankID:         123,
		BaseCurrency:   "USD",
		TargetCurrency: "GBP",
		Tier:           "1",
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
}
