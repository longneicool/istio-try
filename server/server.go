package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"

	istiotest "github.com/longneicool/istio-try"
	"google.golang.org/grpc"
)

type server struct{}

type myHandler struct{}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Rcv message from %s", r.URL.Path)
	fmt.Fprintf(w, "Rcv the request from %s", r.URL.Path)
}

func ListeningHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":9091", nil))
}

func (s *server) SendMessage(ctx context.Context, in *istiotest.Request) (*istiotest.Reply, error) {
	fmt.Printf("Rcv message from %s:%s\n", in.GetName(), in.GetMessage())
	var reply istiotest.Reply
	reply.Status = 1

	fmt.Printf("Send OK to the client %s\n", in.GetName())
	return &reply, nil
}

func main() {
	go ListeningHttp()

	name := os.Getenv("SERVER_NAME")
	if len(name) == 0 {
		name = "NONE"
	}

        ip := os.Getenv("podIP")
	if len(ip) == 0 {
		log.Println("The env podIP is not set, use default one 127.0.0.1")
		ip = "127.0.0.1"
	}

	lis, err := net.Listen("tcp", ip + ":" + "9092")
	if err != nil {
		log.Fatalf("listening failed!err:%v", err)
	}

	grpcServer := grpc.NewServer()
	istiotest.RegisterRoutMessageServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Serve failed!%s", err.Error())
	}

	log.Printf("start listening on %s:%s\n in %s", ip, "9092", name)
}
