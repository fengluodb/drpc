package math

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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

func MarshalMathRequest(v *MathRequest) []byte {
	data, _ := v.Marshal()
	return data
}

func UnmarshalMathRequest(r io.Reader) *MathRequest {
	br := r.(*bytes.Reader)
	data, _ := io.ReadAll(br)
	v := new(MathRequest)
	v.Unmarshal(data)
	tmp := MarshalMathRequest(v)
	*br = *bytes.NewReader(data[len(tmp):])

	return v
}

func MarshalMathReply(v *MathReply) []byte {
	data, _ := v.Marshal()
	return data
}

func UnmarshalMathReply(r io.Reader) *MathReply {
	br := r.(*bytes.Reader)
	data, _ := io.ReadAll(br)
	v := new(MathReply)
	v.Unmarshal(data)
	tmp := MarshalMathReply(v)
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
func (x *MathRequest) Marshal() ([]byte, error) {
	data := []byte{}

	if x.A != 0 {
		data = append(data, MarshalUint8(1)...)
		data = append(data, MarshalInt32(x.A)...)
	} else {
		return nil, fmt.Errorf("marshal failed, A must have value")
	}

	if x.B != 0 {
		data = append(data, MarshalUint8(2)...)
		data = append(data, MarshalInt32(x.B)...)
	} else {
		return nil, fmt.Errorf("marshal failed, B must have value")
	}

	return data, nil
}

func (x *MathReply) Marshal() ([]byte, error) {
	data := []byte{}

	if x.C != 0 {
		data = append(data, MarshalUint8(1)...)
		data = append(data, MarshalInt32(x.C)...)
	} else {
		return nil, fmt.Errorf("marshal failed, C must have value")
	}

	return data, nil
}

func (x *MathRequest) Unmarshal(data []byte) error {
	r := bytes.NewReader(data)

	seq := UnmarshalUint8(r)
	 if seq == 1 {
		x.A = UnmarshalInt32(r)
		seq = UnmarshalUint8(r)
	} else {
		return fmt.Errorf("unmarshal failed, don't find A")
	}

	 if seq == 2 {
		x.B = UnmarshalInt32(r)
	} else {
		return fmt.Errorf("unmarshal failed, don't find B")
	}

	return nil
}

func (x *MathReply) Unmarshal(data []byte) error {
	r := bytes.NewReader(data)

	seq := UnmarshalUint8(r)
	 if seq == 1 {
		x.C = UnmarshalInt32(r)
	} else {
		return fmt.Errorf("unmarshal failed, don't find C")
	}

	return nil
}

