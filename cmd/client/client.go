package client

import (
	"context"
	"fmt"
	"time"

	"github.com/asoorm/todo-grpc/pkg/model/v1/address_formatter"

	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/asoorm/todo-grpc/pkg/model/v1/ping_pong"
	v1 "github.com/asoorm/todo-grpc/pkg/service/v1"
	"google.golang.org/grpc"
)

func Run(grpcServicePort int) {
	serverAddress := fmt.Sprintf("localhost:%d", grpcServicePort)
	log.Info("connecting to grpc server %s...", serverAddress)
	conn, err := grpc.Dial(serverAddress,
		grpc.WithInsecure(),
		//grpc.WithBlock(),
	)
	log.FatalOnError(err)
	log.Info("...connected")
	defer conn.Close()

	pingClient := v1.NewPingPongService()
	addressFormatterClient := v1.NewAddressFormatterService()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for range time.Tick(2 * time.Second) {
		message := fmt.Sprintf("Hello %s", time.Now())
		log.Info("sending: %s", message)
		r, err := pingClient.Ping(ctx, &ping_pong.PingMessage{
			ApiVersion: "v1",
			Message:    message,
		})
		log.FatalOnError(err)
		log.Info("Got %#v", r.GetMessage())

		log.Info("***************************")

		afRes, err := addressFormatterClient.Format(ctx, &address_formatter.AddressRequest{
			ApiVersion: "v1",
			BillingAddress: &address_formatter.Address{
				StreetAddress: "street",
				City:          "city",
				State:         "state",
			},
			ShippingAddress: &address_formatter.Address{
				StreetAddress: "shipping_street",
				City:          "shipping_city",
				State:         "shipping_state",
			},
		})
		log.FatalOnError(err)
		log.Info("AFResponse: Billing %#v", afRes.GetBillingAddress())
		log.Info("AFResponse: Shipping %#v", afRes.GetShippingAddress())

		log.Info("===========================")
	}
}
