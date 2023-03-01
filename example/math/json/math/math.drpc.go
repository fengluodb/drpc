package math

import "github.com/fengluodb/drpc"


type Math interface {
	Add(*MathRequest, *MathReply) error
}

type MathHandler interface {
	AddHandler(req []byte) (data []byte, err error)
}

type MathComplement struct {
	Math Math
}

func (c *MathComplement) AddHandler(req []byte) (data []byte, err error) {
	args := new(MathRequest)
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}
	
	reply := new(MathReply)
	if err := c.Math.Add(args, reply); err != nil {
		return nil, err
	}
	return reply.Marshal()
}

func RegisterMathService(s *drpc.Server, serviceName string, complement Math) {
	c := &MathComplement{
		Math: complement,
	}
	
	drpc.RegisterService(s, serviceName+".Add", c.AddHandler)
}
