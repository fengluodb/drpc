package math

import "encoding/json"


var _ Serializer = (*MathRequest)(nil)
var _ Serializer = (*MathReply)(nil)

type Serializer interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type MathRequest struct {
	A int32
	B int32
}

type MathReply struct {
	C int32
}

func (x *MathRequest) Marshal() ([]byte, error) {
	return json.Marshal(x)
}

func (x *MathRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, x)
}

func (x *MathReply) Marshal() ([]byte, error) {
	return json.Marshal(x)
}

func (x *MathReply) Unmarshal(data []byte) error {
	return json.Unmarshal(data, x)
}
