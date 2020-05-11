# REST proxy fascade for gRPC services

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

## Start gRPC server (with server reflection enabled)

```shell script

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

## Notes:

Log of notes, & reading material & resources for this project

[Scrapbook](/SCRAPBOOK.md)
