package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/govardhanpagidi/nats-client/fxconversion"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect("nats://localhost:4222")
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
	} // Convert JSON to Protocol Buffer message

	// Marshal your Protocol Buffer message
	protoBytes, err := proto.Marshal(yourMessage)
	if err != nil {
		log.Fatal(err)
	}

	// Your gRPC message payload (example byte array)
	//grpcMessage := []byte("Your gRPC message")

	// Publish gRPC message to NATS subject
	err = nc.Publish(subject, protoBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Published gRPC message to NATS subject:", subject)
}
