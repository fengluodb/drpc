package drpc

import (
	"encoding/binary"
)

const (
	// MaxHeaderSize = 10 + 10 + 4 (10 refer to binary.MaxVarintLen64)
	MaxHeaderSize = 24

	Uint16Size = 2
	Uint32Size = 4
)

// RequestHeader request header structure looks like:
// +----------+----------------+----------+
// |    ID    |      Method    |  Checksum|
// +----------+----------------+----------+
// |  uvarint | uvarint+string |   uint32 |
// +----------+----------------+----------+
type RequestHeader struct {
	ID       uint64
	Method   string
	Checksum uint32
}

func (r *RequestHeader) Marshal() []byte {
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Method))

	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += writeString(header[idx:], r.Method)
	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size

	return header[:idx]
}

func (r *RequestHeader) Unmarshal(data []byte) error {
	idx, size := 0, 0
	n := len(data)

	if idx >= n {
		return ErrUnmarshal
	}
	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	if idx >= n {
		return ErrUnmarshal
	}
	r.Method, size = readString(data[idx:])
	idx += size

	if idx >= n {
		return ErrUnmarshal
	}
	r.Checksum = binary.LittleEndian.Uint32(data[idx:])

	return nil
}

// ResponseHeader request header structure looks like:
// +---------+----------------+----------+
// |    ID   |      Error     |  Checksum|
// +---------+----------------+----------+
// | uvarint | uvarint+string |   uint32 |
// +---------+----------------+----------+
type ResponseHeader struct {
	ID       uint64
	Error    string
	Checksum uint32
}

func (r *ResponseHeader) Marshal() []byte {
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Error))

	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += writeString(header[idx:], r.Error)
	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size

	return header[:idx]
}

func (r *ResponseHeader) Unmarshal(data []byte) error {
	idx, size := 0, 0
	n := len(data)

	if idx >= n {
		return ErrUnmarshal
	}
	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	if idx >= n {
		return ErrUnmarshal
	}
	r.Error, size = readString(data[idx:])
	idx += size

	if idx >= n {
		return ErrUnmarshal
	}
	r.Checksum = binary.LittleEndian.Uint32(data[idx:])

	return nil
}

func readString(data []byte) (string, int) {
	idx := 0
	length, size := binary.Uvarint(data)
	idx += size
	str := string(data[idx : idx+int(length)])
	idx += len(str)
	return str, idx
}

func writeString(data []byte, str string) int {
	idx := 0
	idx += binary.PutUvarint(data, uint64(len(str)))
	copy(data[idx:], str)
	idx += len(str)
	return idx
}
