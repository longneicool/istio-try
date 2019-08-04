package main

import (
	"context"
	"flag"
	"log"
	"time"

	istiotest "github.com/longneicool/istio-try"
	"google.golang.org/grpc"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "the ip of the server")
	port := flag.String("port", "50001", "the port of the server")
	name := flag.String("name", "client1", "The name of the client")
	flag.Parse()

	conn, err := grpc.Dial(*ip+":"+*port, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := istiotest.NewRoutMessageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := istiotest.Request{
		Name:    *name,
		Message: "I am " + *name,
	}

	reply, err := client.SendMessage(ctx, &req)
	if err != nil {
		log.Fatal("Rcv error from the server!")
	}

	log.Println("Rcv response ", reply.GetStatus())
}
