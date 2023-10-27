package methodhandler

import (
	"google.golang.org/grpc"
)

type Handler interface {
	Name() string
	Handle(stream grpc.ServerStream) error
}
