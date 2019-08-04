package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	istiotest "github.com/longneicool/istio-try"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) SendMessage(ctx context.Context, in *istiotest.Request) (*istiotest.Reply, error) {
	fmt.Printf("Rcv message from %s:%s\n", in.GetName(), in.GetMessage())
	var reply istiotest.Reply
	reply.Status = 1

	fmt.Printf("Send OK to the client %s\n", in.GetName())
	return &reply, nil
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "The listener ip of the server")
	port := flag.String("prot", "50001", "the listener port")
	flag.Parse()
	fmt.Printf("Listening on %s:%s\n", *ip, *port)

	lis, err := net.Listen("tcp", *ip+":"+*port)
	if err != nil {
		log.Fatal("listening failed!err:%v", err)
	}

	grpcServer := grpc.NewServer()
	istiotest.RegisterRoutMessageServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Serve failed!%s", err.Error())
	}
}
