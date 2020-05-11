# REST proxy facade for gRPC services

## Installation

1. Generate Self-Signed Certs (plaintext is currently not enabled)

```shell script
sh certs.sh

Generating RSA private key, 2048 bit long modulus
..........................+++
..............................+++
e is 65537 (0x10001)
```

This creates `server.crt` and `server.key` in the `./certs` directory.

2. Start gRPC server (with server reflection enabled)

```shell script
go run main.go

INFO: starting grpc server
INFO: os.Getwd(): /Users/ahmet/go/src/github.com/asoorm/grpc-rest
INFO: os.Executable(): /var/folders/73/x7scyhbx0cj5llzdc9k99jzr0000gn/T/go-build816017477/b001/exe/main
INFO: listening on 9000
```

You can also be explicit to start the server & override the listen port

```shell script
go run main.go -mode server -port 9191

INFO: starting grpc server
INFO: os.Getwd(): /Users/ahmet/go/src/github.com/asoorm/grpc-rest
INFO: os.Executable(): /var/folders/73/x7scyhbx0cj5llzdc9k99jzr0000gn/T/go-build094596563/b001/exe/main
INFO: listening on 9191
```

If you wish to enable transport logging verbosity for debug purposes
```shell script
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run main.go

INFO: starting grpc server
INFO: os.Getwd(): /Users/ahmet/go/src/github.com/asoorm/grpc-rest
INFO: os.Executable(): /var/folders/73/x7scyhbx0cj5llzdc9k99jzr0000gn/T/go-build444463768/b001/exe/main
INFO: listening on 9000
WARNING: 2020/05/11 07:42:17 grpc: Server.Serve failed to complete security handshake from "[::1]:55158": tls: first record does not look like a TLS handshake
WARNING: 2020/05/11 07:42:18 grpc: Server.Serve failed to complete security handshake from "[::1]:55170": tls: first record does not look like a TLS handshake
WARNING: 2020/05/11 07:42:20 grpc: Server.Serve failed to complete security handshake from "[::1]:55171": tls: first record does not look like a TLS handshake
WARNING: 2020/05/11 07:42:23 grpc: Server.Serve failed to complete security handshake from "[::1]:55172": tls: first record does not look like a TLS handshake
WARNING: 2020/05/11 07:42:28 grpc: Server.Serve failed to complete security handshake from "[::1]:55173": tls: first record does not look like a TLS handshake
```

3. Start the REST Proxy (PoC + Buggy)

```shell script
go run main.go --mode rest
   
INFO: starting rest proxy
INFO: reflected services: []string{"grpc.reflection.v1alpha.ServerReflection", "v1.PingPongService"}
INFO: service 0: ServerReflection
INFO: /ServerReflection
INFO:   /ServerReflection/ServerReflectionInfo
INFO: serviceDescriptor name:"ServerReflection" method:{name:"ServerReflectionInfo" input_type:".grpc.reflection.v1alpha.ServerReflectionRequest" output_type:".grpc.reflection.v1alpha.ServerReflectionResponse" client_streaming:true server_streaming:true}
INFO: service 1: PingPongService
INFO: /PingPongService
INFO:   /PingPongService/Ping
INFO: serviceDescriptor name:"PingPongService" method:{name:"Ping" input_type:".v1.PingMessage" output_type:".v1.PongMessage" options:{}}
INFO: starting rest server on :9001
```

Using server reflection, we have identified a couple of services. The `PingPongService` & the `Ping` method.

The REST proxy creates the router based on this information:

```shell script
/PingPongService/Ping
```

The proxy is able to dynamically marshal & unmarshal the JSON request/response bodies to the appropriate
input and output types.

4. Send REST calls to the proxy

Invalid
```shell script
curl http://localhost:9001/PingPongService/Ping -d '{"api_version": "v2", "message": "hello ahmet"}'

Resolved method descriptor:
rpc Ping ( .v1.PingMessage ) returns ( .v1.PongMessage );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
```

Unknown path

```shell script
curl http://localhost:9001/PingPongService/Pong -d '{"api_version": "v1", "message": "hello ahmet"}'

404 page not found
```

Valid
```shell script
curl http://localhost:9001/PingPongService/Ping -d '{"api_version": "v1", "message": "hello ahmet"}'

Resolved method descriptor:
rpc Ping ( .v1.PingMessage ) returns ( .v1.PongMessage );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc

Response contents:
{
  "apiVersion": "v1",
  "message": "PONG: hello ahmet"
}

Response trailers received:
(empty)
```

## Notes:

Log of notes, & reading material & resources for this project

[Scrapbook](/SCRAPBOOK.md)
