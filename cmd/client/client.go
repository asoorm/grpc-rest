package client

import (
	"context"
	"fmt"
	"time"

	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/asoorm/todo-grpc/pkg/model/v1/ping_pong"
	v1 "github.com/asoorm/todo-grpc/pkg/service/v1"
	"google.golang.org/grpc"
)

func Run(port int) {

	serverAddress := fmt.Sprintf("localhost:%d", port)
	log.Info("connecting to grpc server %s...", serverAddress)
	conn, err := grpc.Dial(serverAddress,
		grpc.WithInsecure(),
		//grpc.WithBlock(),
	)
	log.FatalOnError(err)
	log.Info("...connected")
	defer conn.Close()

	c := v1.NewPingPongService()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for range time.Tick(2 * time.Second) {
		message := fmt.Sprintf("Hello %s", time.Now())
		log.Info("sending: %s", message)
		r, err := c.Ping(ctx, &ping_pong.PingMessage{
			ApiVersion: "v1",
			Message:    message,
		})
		log.FatalOnError(err)
		log.Info("Got %#v", r.GetMessage())
	}
}
