# Description

This `protoc-gen-gofullmethods` plugin is intended as an extension of the [protoc-gen-go-grpc plugin](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc), exposing the full methods for each gRPC endpoint as string constants. Its original intent is to use these methods/paths for authorization purposes in a custom middleware function in the transport layer, on the gRPC server. This plugin is not intended for use with clients.

The `example` directory uses the sample pet API from buf to generate example output of this plugin, using the locally compiled executable. Run `make update-example` to regenerate it.

## Requirements

- [golang](https://golang.org/) v1.19+

## Other plugin dependencies

As per `go.mod`:

- google.golang.org/protobuf:v1.27.0

## Installing

```bash
go install github.com/90poe/protoc-gen-gofullmethods@latest
```
