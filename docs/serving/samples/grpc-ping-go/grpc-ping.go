package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	ping "github.com/knative/docs/docs/serving/samples/grpc-ping-go/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	value = flag.String("value", "pong", "The string to return.")
)

type pingServer struct {
	Value string
}

func (p *pingServer) Ping(ctx context.Context, req *ping.Request) (*ping.Response, error) {
	return &ping.Response{Msg: fmt.Sprintf("%s - %s", req.Msg, p.Value)}, nil
}

func (p *pingServer) PingStream(stream ping.PingService_PingStreamServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("Client disconnected")
			return nil
		}

		if err != nil {
			fmt.Println("Failed to receive ping")
			return err
		}

		fmt.Printf("Replying to ping %s at %s\n", req.Msg, time.Now())

		err = stream.Send(&ping.Response{
			Msg: fmt.Sprintf("%s %s", p.Value, time.Now()),
		})

		if err != nil {
			fmt.Printf("Failed to send pong %s\n", err)
			return err
		}
	}
}

func main() {
	flag.Parse()

	// Default the port to 8080 if unspecified.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pingServer := &pingServer{
		Value: *value,
	}

	// The grpcServer is currently configured to serve h2c traffic by default.
	// To configure credentials or encryption, see: https://grpc.io/docs/guides/auth.html#go
	grpcServer := grpc.NewServer()
	ping.RegisterPingServiceServer(grpcServer, pingServer)
	grpcServer.Serve(lis)
}
