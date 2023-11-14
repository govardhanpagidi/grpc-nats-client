package grpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/govardhanpagidi/nats-client/proto"
	"github.com/nats-io/nats.go"
	"log"
)

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
