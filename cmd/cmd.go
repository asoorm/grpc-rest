package cmd

import (
	"flag"
	"fmt"

	"github.com/asoorm/todo-grpc/pkg/log"

	"github.com/asoorm/todo-grpc/cmd/server"
)

var (
	mode *string
	port *int
)

func Run() {
	mode = flag.String("mode", "server", "start as client or server {client,server,rest}[default:server]")
	port = flag.Int("port", 9000, "grpc server port [default:9000]")
	flag.Parse()

	switch *mode {
	case "server":
		log.Info("starting grpc server")
		server.Run(*port)
	//case "client":
	//	log.Info("starting grpc client")
	//	client.Run(*port)
	//case "rest":
	//	log.Info("starting rest proxy")
	//	rest_proxy.Run(*port)
	default:
		log.FatalOnError(fmt.Errorf("unsupported mode, expected {client,server}, got %s", *mode))
	}
}
