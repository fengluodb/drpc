package drpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mathArgs struct {
	A int
	B int
}

type mathReply struct {
	C int
}

func (ma *mathArgs) Marshal() ([]byte, error) {
	return json.Marshal(ma)
}

func (ma *mathArgs) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ma)
}

func (mr *mathReply) Marshal() ([]byte, error) {
	return json.Marshal(mr)
}

func (mr *mathReply) Unmarshal(data []byte) error {
	return json.Unmarshal(data, mr)
}

type math struct{}

func (m *math) Add(args *mathArgs, reply *mathReply) error {
	reply.C = args.A + args.B
	return nil
}

func (m *math) Mul(args *mathArgs, reply *mathReply) error {
	reply.C = args.A * args.B
	return nil
}

// dgen generate
type Math interface {
	Add(*mathArgs, *mathReply) error
	Mul(*mathArgs, *mathReply) error
}

var _ MathHandler = (*mathHandler)(nil)

type MathHandler interface {
	AddHandler(req []byte) (data []byte, err error)
	MulHandler(req []byte) (data []byte, err error)
}

type mathHandler struct {
	math Math
}

func (m *mathHandler) AddHandler(req []byte) (data []byte, err error) {
	args := new(mathArgs)
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}

	reply := new(mathReply)

	if err := m.math.Add(args, reply); err != nil {
		return nil, err
	}
	return reply.Marshal()
}

func (m *mathHandler) MulHandler(req []byte) (data []byte, err error) {
	args := new(mathArgs)
	if err := args.Unmarshal(req); err != nil {
		return nil, err
	}

	reply := new(mathReply)

	if err := m.math.Mul(args, reply); err != nil {
		return nil, err
	}
	return reply.Marshal()
}

func RegisterMethodService(s *Server, serviceName string, math Math) {
	mathHandler := &mathHandler{
		math: math,
	}
	RegisterService(s, serviceName+".Add", mathHandler.AddHandler)
	RegisterService(s, serviceName+".Mul", mathHandler.MulHandler)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	startServer()
}

func startServer() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Printf("failed to start server, error: %v", err)
	}
	server := NewServer()
	RegisterMethodService(server, "Math", new(math))

	go server.Serve(listener)
}

func TestMath(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(conn)
	for i := 0; i < 100; i++ {
		args := &mathArgs{
			A: rand.Int(),
			B: rand.Int(),
		}
		reply := new(mathReply)

		err = client.Call("Math.Add", args, reply)
		assert.NoError(t, err)
		assert.Equal(t, args.A+args.B, reply.C)

		err = client.Call("Math.Mul", args, reply)
		assert.NoError(t, err)
		assert.Equal(t, args.A*args.B, reply.C)
	}
}

func TestConcurrentMath(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(conn)
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("g-%d", i), func(t *testing.T) {
			t.Parallel()

			args := &mathArgs{
				A: rand.Int(),
				B: rand.Int(),
			}
			reply := new(mathReply)

			err = client.Call("Math.Add", args, reply)
			assert.NoError(t, err)
			assert.Equal(t, args.A+args.B, reply.C)

			err = client.Call("Math.Mul", args, reply)
			assert.NoError(t, err)
			assert.Equal(t, args.A*args.B, reply.C)
		})
	}
}
