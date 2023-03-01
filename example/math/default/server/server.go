package main

import (
	"fmt"
	"net"

	"github.com/fengluodb/drpc"
	dgen "github.com/fengluodb/drpc/example/math/default/math"
)

type mathService struct{}

func (m *mathService) Add(args *dgen.MathRequest, reply *dgen.MathReply) error {
	reply.C = args.A + args.B
	return nil
}

func main() {
	server := drpc.NewServer()
	dgen.RegisterMathService(server, "Math", new(mathService))

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}
	server.Serve(listener)
}
