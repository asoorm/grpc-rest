package cmd

import (
	"flag"
	"fmt"

	"github.com/asoorm/todo-grpc/cmd/client"
	"github.com/asoorm/todo-grpc/cmd/rest_proxy"
	"github.com/asoorm/todo-grpc/cmd/server"
	"github.com/asoorm/todo-grpc/pkg/log"
)

var (
	mode     = flag.String("mode", "server", "start as client or server {client,server,rest}[default:server]")
	addr     = flag.String("addr", "localhost:9000", "gRPC server port [localhost:9000]")
	restPort = flag.Int("rest_port", 9001, "REST API port [default:9001], only valid with {mode:rest}")
)

func Run() {
	flag.Parse()

	switch *mode {
	case "server":
		log.Info("starting grpc server")
		server.Run(*addr)
	case "client":
		log.Info("starting grpc client")
		client.Run(*addr)
	case "rest":
		log.Info("starting rest proxy")
		rest_proxy.Run(*restPort, *addr)
	default:
		log.FatalOnError(fmt.Errorf("unsupported mode, expected {client,server,rest}, got %s", *mode))
	}
}
