package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc/credentials"
	"time"

	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/asoorm/todo-grpc/pkg/model/v1/ping_pong"
	"google.golang.org/grpc"
)

func Run(addr string) {
	log.Info("connecting to grpc server %s...", addr)
	conn, err := grpc.Dial(addr,
		//grpc.WithInsecure(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
		grpc.WithBlock(),
	)
	log.FatalOnError(err)
	log.Info("...connected")
	defer conn.Close()

	pingClient := ping_pong.NewPingPongServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

		//afRes, err := addressFormatterClient.Format(ctx, &address_formatter.AddressRequest{
		//	ApiVersion: "v1",
		//	BillingAddress: &address_formatter.Address{
		//		StreetAddress: "street",
		//		City:          "city",
		//		State:         "state",
		//	},
		//	ShippingAddress: &address_formatter.Address{
		//		StreetAddress: "shipping_street",
		//		City:          "shipping_city",
		//		State:         "shipping_state",
		//	},
		//})
		//log.FatalOnError(err)
		//log.Info("AFResponse: Billing %#v", afRes.GetBillingAddress())
		//log.Info("AFResponse: Shipping %#v", afRes.GetShippingAddress())

		log.Info("===========================")
	}
}
