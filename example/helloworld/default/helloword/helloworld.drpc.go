package helloworld

import "github.com/fengluodb/drpc"

type HelloWorld interface {
	SayHello(*HelloRequest, *HelloReply) error
}

type HelloWorldHandler interface {
	SayHelloHandler(req []byte) (data []byte, err error)
}

type HelloWorldComplement struct {
	HelloWorld HelloWorld
}

func (c *HelloWorldComplement) SayHelloHandler(req []byte) (data []byte, err error) {
	args := new(HelloRequest)
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}

	reply := new(HelloReply)
	if err := c.HelloWorld.SayHello(args, reply); err != nil {
		return nil, err
	}
	return reply.Marshal()
}

func RegisterHelloWorldService(s *drpc.Server, serviceName string, complement HelloWorld) {
	c := &HelloWorldComplement{
		HelloWorld: complement,
	}

	drpc.RegisterService(s, serviceName+".SayHello", c.SayHelloHandler)
}
