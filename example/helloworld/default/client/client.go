package main

import (
	"fmt"

	"github.com/fengluodb/drpc"
	dgen "github.com/fengluodb/drpc/example/helloworld/default/helloword"
)

func main() {
	client, err := drpc.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}

	args := &dgen.HelloRequest{
		Name: "fengluodb",
	}
	reply := &dgen.HelloReply{}

	if err := client.Call("Hello.SayHello", args, reply); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("reply:", reply.Reply)
}
