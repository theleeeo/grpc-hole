# gRPC-Hole

gRPC-Hole is a tool for mocking gRPC servers.

## Terminology
- Service: A protobuf service. It is the api-definition of what methods exists and what types they use.
- Server: A gRPC-Hole, simply called a server, is a program listening for input requests and creates a response based. A server is based on a Service to specify which methods to handle.
- Method: An RPC (Remote Procedure Call) that can be called. They are the API endpoints of a server.

## Services

gRPC-Hole uses defined protobuf services to as a base for a server. They contain which methods exists, what input/output types they have and what extentions are present.

To list all currently saved services:

```
grpc-hole services list
```

### Scanning

Through the grpc-hole cli you can scan a running gRPC-server using reflection to copy its services.

The server must have reflection enabled to be able to be scanned.

How to enable reflection:
- [Go](https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md)
- [Java](https://github.com/grpc/grpc-java/blob/master/documentation/server-reflection-tutorial.md)
- [C++](https://github.com/grpc/grpc/blob/master/doc/server_reflection_tutorial.md)
- [Python](https://github.com/grpc/grpc/blob/master/doc/python/server_reflection.md)

Reflection is supported for more languages than those but that is left as an exercise for the reader.

Example of how to use the cli to scan a service:

```
grpc-hole services scan -t localhost:5000
```

### Loading

Loading a service from proto-descriptors is coming soon!

## Servers

### Static Server

A static server is responding to the requests with a pre-defined response. The static response can however include go-templates to vary the response based on the input.

Learn more about go-templates [here](https://golang.org/pkg/text/template/)

### Proxy Server

A proxy server is a server that forwards the requests to another server. It can be used to log the requests and responses that are being passed through.

Coming soon:
- Modify requests before forwarding
- Modify responses before returning them
