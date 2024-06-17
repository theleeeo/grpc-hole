FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-hole .

FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/grpc-hole .

CMD ["./grpc-hole"]
