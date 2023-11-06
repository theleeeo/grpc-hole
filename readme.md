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

A service can also be loaded from a .proto file. This is useful if you want to create a server from a service that is not running or if a server does not have reflection enebled but you have the definition.

Example of how to use the cli to load a service:

```
grpc-hole services load -r /path/to/root -f /relative/to/root.proto
```

The flag `-r` is the root directory of the protos. This is used as the base to resolve imports and the main files containing the services.

The flag `-f` is the path to the file/files containing the services. This is relative to the root directory specified with `-r`.

#### Known issues:
If a service is loaded from a file (or with a dependency of a file) that defines a custom option and that option is used within the same file it is defined in, the option will not be linked correctly and the service will be broken. This is due to the way the protoc compiler works and is not a bug in gRPC-Hole.

To ensure this will not be an issue, make sure to define all custom options in a separate file and import it into the file where it is used.

## Servers

### Static Server

A static server is responding to the requests with a pre-defined response. The static response can however include go-templates to vary the response based on the input.

Learn more about go-templates [here](https://golang.org/pkg/text/template/)

A static server can be started with the following command:


```
grpc-hole server static -s=ServiceName -p=Port
```

The flag `-s` is the name of the service to use and the `-p` is the port to listen on.

If the host must be specified, use the `--host` flag. The default host is `0.0.0.0`.

### Proxy Server

A proxy server is a server that forwards the requests to another server. It can be used to log the requests and responses that are being passed through.

A proxy server can be started with the following command:

```
grpc-hole server proxy -s=ServiceName -p=Port -t=Target
```

The flag `-s` is the name of the service to use and the `-p` is the port to listen on.

The flag `-t` is the target address to forward the requests to.

Coming soon:
- Modify requests before forwarding
- Modify responses before returning them
