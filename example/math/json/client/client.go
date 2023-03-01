package main

import (
	"fmt"
	"time"

	"github.com/fengluodb/drpc"
	dgen "github.com/fengluodb/drpc/example/math/json/math"
)

func main() {
	client, err := drpc.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}

	args := &dgen.MathRequest{
		A: 1,
		B: 2,
	}
	reply := &dgen.MathReply{}

	t := time.Now()
	for i := 0; i < 100; i++ {
		for j := 0; j < 1000; j++ {
			if err := client.Call("Math.Add", args, reply); err != nil {
				fmt.Println("err:", err)
			}
		}
	}
	fmt.Println(time.Now().Sub(t).Seconds() / 100)
}
