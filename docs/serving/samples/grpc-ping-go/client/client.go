package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/knative/docs/docs/serving/samples/grpc-ping-go/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	server   = flag.String("server", "127.0.0.1", "The server hostname.")
	insecure = flag.Bool("insecure", false, "Set to true to skip SSL validation")
)

func main() {
	flag.Parse()

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Unable to read system certificates")
	}
	creds := credentials.NewClientTLSFromCert(pool, "")

	opts := []grpc.DialOption{}
	var port int32
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
		port = 80
	} else {
		opts = append(opts, grpc.WithTransportCredentials(creds))
		port = 443
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *server, port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)

	ping(client, "hello")
	pingStream(client, "hello")
}

func ping(client pb.PingServiceClient, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	rep, err := client.Ping(ctx, &pb.Request{Msg: msg})
	if err != nil {
		log.Fatalf("%v.Ping failed %v: ", client, err)
	}
	log.Printf("Ping got %v\n", rep.GetMsg())
}

func pingStream(client pb.PingServiceClient, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.PingStream(ctx)
	if err != nil {
		log.Fatalf("%v.(_) = _, %v", client, err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a response : %v", err)
			}
			log.Printf("Got %s", in.GetMsg())
		}
	}()

	i := 0
	for i < 20 {
		if err := stream.Send(&pb.Request{Msg: fmt.Sprintf("%s-%d", msg, i)}); err != nil {
			log.Fatalf("Failed to send a ping: %v", err)
		}
		i++
	}
	stream.CloseSend()
	<-waitc

}
