package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
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

	conversionReqJsonStr := `{
		"tenantId": 1235,
		"bankId": 123456,
		"baseCurrency": "GBP",
		"targetCurrency": "USD",
		"tier": "2",
		"amount":1000
	}`
	yourMessage := &pb.ConversionRequest{} // Convert JSON to Protocol Buffer message
	err = json.Unmarshal([]byte(conversionReqJsonStr), yourMessage)
	if err != nil {
		log.Fatal(err)
	}

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
