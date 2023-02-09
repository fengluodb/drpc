package main

import (
	"fmt"
	"net"

	"github.com/fengluodb/drpc"
	dgen "github.com/fengluodb/drpc/example/helloworld/default/helloword"
)

type helloService struct {
}

func (h *helloService) SayHello(args *dgen.HelloRequest, reply *dgen.HelloReply) error {
	reply.Reply = fmt.Sprintf("Hello %s, welcome to drpc", args.Name)
	return nil
}

func main() {
	server := drpc.NewServer()
	dgen.RegisterHelloWorldService(server, "Hello", new(helloService))

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}
	server.Serve(listener)
}
