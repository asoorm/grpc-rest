package server

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/asoorm/todo-grpc/pkg/model/v1/ping_pong"
	pb "github.com/asoorm/todo-grpc/pkg/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func Run(port int) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	log.FatalOnError(err)

	workingDir, _ := os.Getwd()
	executable, _ := os.Executable()

	log.Info("os.Getwd(): %s", workingDir)
	log.Info("os.Executable(): %s", executable)

	creds, err := credentials.NewServerTLSFromFile(workingDir+"/certs/server.crt", workingDir+"/certs/server.key")
	log.FatalOnError(err)

	pingPongService := pb.NewPingPongService()
	//addressFormatterService := pb.NewAddressFormatterService()
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	ping_pong.RegisterPingPongServiceServer(grpcServer, pingPongService)
	//address_formatter.RegisterAddressFormatterServiceServer(grpcServer, addressFormatterService)
	reflection.Register(grpcServer)

	log.Info("listening on %d", port)
	log.FatalOnError(grpcServer.Serve(lis))
}
