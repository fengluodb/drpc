package drpc

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestRequestHeader(t *testing.T) {
	for i := 0; i < 100; i++ {
		src := GenerateRandomRequestHeader()
		data := src.Marshal()
		dest := new(RequestHeader)
		err := dest.Unmarshal(data)

		assert.NoError(t, err)
		assert.Equal(t, src, dest)
	}
}

func TestResponseHeader(t *testing.T) {
	for i := 0; i < 100; i++ {
		src := GenerateRandomResponseHeader()
		data := src.Marshal()
		dest := new(ResponseHeader)
		err := dest.Unmarshal(data)

		assert.NoError(t, err)
		assert.Equal(t, src, dest)
	}
}

func GenerateRandomRequestHeader() *RequestHeader {
	return &RequestHeader{
		ID:       rand.Uint64(),
		Method:   GetRandomString(),
		Checksum: rand.Uint32(),
	}
}

func GenerateRandomResponseHeader() *ResponseHeader {
	return &ResponseHeader{
		ID:       rand.Uint64(),
		Error:    GetRandomString(),
		Checksum: rand.Uint32(),
	}
}

func GetRandomString() string {
	randBytes := make([]byte, rand.Intn(1000))
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
