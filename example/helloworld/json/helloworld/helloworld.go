package helloworld

import "encoding/json"


type color uint32

var _ Serializer = (*HelloRequest)(nil)
var _ Serializer = (*HelloReply)(nil)

type Serializer interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type HelloRequest struct {
	Name string
}

type HelloReply struct {
	Reply string
}

func (x *HelloRequest) Marshal() ([]byte, error) {
	return json.Marshal(x)
}

func (x *HelloRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, x)
}

func (x *HelloReply) Marshal() ([]byte, error) {
	return json.Marshal(x)
}

func (x *HelloReply) Unmarshal(data []byte) error {
	return json.Unmarshal(data, x)
}
