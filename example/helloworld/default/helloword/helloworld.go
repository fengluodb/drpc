package helloworld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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

func MarshalHelloRequest(v *HelloRequest) []byte {
	data, _ := v.Marshal()
	return data
}

func UnmarshalHelloRequest(r io.Reader) *HelloRequest {
	br := r.(*bytes.Reader)
	data, _ := io.ReadAll(br)
	v := new(HelloRequest)
	v.Unmarshal(data)
	tmp := MarshalHelloRequest(v)
	*br = *bytes.NewReader(data[len(tmp):])

	return v
}

func MarshalHelloReply(v *HelloReply) []byte {
	data, _ := v.Marshal()
	return data
}

func UnmarshalHelloReply(r io.Reader) *HelloReply {
	br := r.(*bytes.Reader)
	data, _ := io.ReadAll(br)
	v := new(HelloReply)
	v.Unmarshal(data)
	tmp := MarshalHelloReply(v)
	*br = *bytes.NewReader(data[len(tmp):])

	return v
}

func MarshalUint8(v uint8) []byte {
	data := []byte{}
	return append(data, byte(v))
}

func UnmarshalUint8(r io.Reader) uint8 {
	data := make([]byte, 1)
	io.ReadFull(r, data)
	return uint8(data[0])
}

func MarshalUint16(v uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, v)
	return data
}

func UnmarshalUint16(r io.Reader) uint16 {
	data := make([]byte, 2)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint16(data)
}

func MarshalUint32(v uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, v)
	return data
}

func UnmarshalUint32(r io.Reader) uint32 {
	data := make([]byte, 4)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint32(data)
}

func MarshalUint64(v uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, v)
	return data
}

func UnmarshalUint64(r io.Reader) uint64 {
	data := make([]byte, 8)
	io.ReadFull(r, data)
	return binary.LittleEndian.Uint64(data)
}

func MarshalInt8(v int8) []byte {
	data := []byte{}
	return append(data, byte(v))
}

func UnmarshalInt8(r io.Reader) int8 {
	data := make([]byte, 1)
	io.ReadFull(r, data)
	return int8(data[0])
}

func MarshalInt16(v int16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(v))
	return data
}

func UnmarshalInt16(r io.Reader) int16 {
	data := make([]byte, 2)
	io.ReadFull(r, data)
	return int16(binary.LittleEndian.Uint16(data))
}

func MarshalInt32(v int32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	return data
}

func UnmarshalInt32(r io.Reader) int32 {
	data := make([]byte, 4)
	io.ReadFull(r, data)
	return int32(binary.LittleEndian.Uint32(data))
}

func MarshalInt64(v int64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(v))
	return data
}

func UnmarshalInt64(r io.Reader) int64 {
	data := make([]byte, 8)
	io.ReadFull(r, data)
	return int64(binary.LittleEndian.Uint64(data))
}

func MarshalString(s string) []byte {
	data := []byte{}
	data = append(data, MarshalInt32(int32(len(s)))...)
	data = append(data, []byte(s)...)
	return data
}

func UnmarshalString(r io.Reader) string {
	size := UnmarshalInt32(r)

	data := make([]byte, size)
	io.ReadFull(r, data)
	return string(data)
}
func (x *HelloRequest) Marshal() ([]byte, error) {
	data := []byte{}

	if x.Name != "" {
		data = append(data, MarshalUint8(1)...)
		data = append(data, MarshalString(x.Name)...)
	} else {
		return nil, fmt.Errorf("marshal failed, Name must have value")
	}

	return data, nil
}

func (x *HelloReply) Marshal() ([]byte, error) {
	data := []byte{}

	if x.Reply != "" {
		data = append(data, MarshalUint8(1)...)
		data = append(data, MarshalString(x.Reply)...)
	}

	return data, nil
}

func (x *HelloRequest) Unmarshal(data []byte) error {
	r := bytes.NewReader(data)

	seq := UnmarshalUint8(r)
	 if seq == 1 {
		x.Name = UnmarshalString(r)
	} else {
		return fmt.Errorf("unmarshal failed, don't find Name")
	}

	return nil
}

func (x *HelloReply) Unmarshal(data []byte) error {
	r := bytes.NewReader(data)

	seq := UnmarshalUint8(r)
	 if seq == 1 {
		x.Reply = UnmarshalString(r)
	}

	return nil
}

