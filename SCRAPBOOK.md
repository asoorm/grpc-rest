# Scrapbook

Notes & learning material used for this project

```
$ grpcurl -insecure localhost:9000 list
grpc.reflection.v1alpha.ServerReflection
v1.PingPongService
```

```
$ grpcurl -insecure localhost:9000 list v1.PingPongService
v1.PingPongService.Ping
```

```
$ grpcurl -insecure localhost:9000 describe v1.PingPongService
v1.PingPongService is a service:
service PingPongService {
  rpc Ping ( .v1.PingMessage ) returns ( .v1.PongMessage );
}
```

```
$ grpcurl -insecure localhost:9000 describe v1.PingMessage    
v1.PingMessage is a message:
message PingMessage {
  string api_version = 1;
  string message = 2;
}
```

```
$ grpcurl -d '{"api_version": "v2", "message": "hi"}' -insecure localhost:9000 v1.PingPongService/Ping
ERROR:
  Code: Unimplemented
  Message: unsupported API requestedVersion, service implements (v1) but requested (v2)
```

```
$ grpcurl -d '{"api_version": "v1", "message": "hi"}' -insecure localhost:9000 v1.PingPongService/Ping
{
  "apiVersion": "v1",
  "message": "PONG: hi"
}
```

---

## gRPC client/server logging

```shell script
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run main.go
```

## Reference

- [grpc.io Docs: Quickstart Go](https://grpc.io/docs/quickstart/go/)
- [github.com/jhump/protoreflect: Protocol Buffer and gRPC Reflection Library](https://github.com/jhump/protoreflect)
- [github.com/protocolbuffers/protobuf-go: 1.20.0 Release Notes](https://github.com/protocolbuffers/protobuf-go/releases/tag/v1.20.0)
- [developers.google.com: Protocol Buffers Go Reference](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
- [bitbucket.org/blog: Writing a microservice in Golang which communicates over gRPC](https://bitbucket.org/blog/writing-a-microservice-in-golang-which-communicates-over-grpc)
- [blog.golang.org/protobuf-apiv2: A new Go API for Protocol Buffers](https://blog.golang.org/protobuf-apiv2)
- [grpc.io/blog/grpc-web-ga: gRPC-Web is Generally Available](https://grpc.io/blog/grpc-web-ga/)
- [medium.com/@amsokol.com: Go gRPC tutorial](https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9)
- [github.com/golang-standards/project-layout: Standard Go Project Layout](https://github.com/golang-standards/project-layout)
