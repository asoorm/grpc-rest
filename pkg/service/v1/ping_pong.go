package v1

import (
	"context"
	"fmt"

	"github.com/asoorm/todo-grpc/pkg/model/v1/ping_pong"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

func NewPingPongService() *PingPongService {
	return &PingPongService{}
}

type PingPongService struct{}

func (s PingPongService) Ping(ctx context.Context, in *ping_pong.PingMessage) (*ping_pong.PongMessage, error) {
	if err := checkAPIVersion(in.ApiVersion); err != nil {
		return nil, err
	}

	return &ping_pong.PongMessage{
		ApiVersion: apiVersion,
		Message:    fmt.Sprintf("PONG: %s", in.GetMessage()),
	}, nil
}

func checkAPIVersion(requestedVersion string) error {
	// requested version supplied? "" means use current version
	if len(requestedVersion) > 0 {
		if requestedVersion != apiVersion {
			return status.Errorf(codes.Unimplemented, "unsupported API requestedVersion, service implements (%s) but requested (%s)", apiVersion, requestedVersion)
		}
	}

	return nil
}
