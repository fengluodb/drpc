package drpc

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"sync"
)

var _ ClientCodec = (*clientCodec)(nil)

type Serializer interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type Call struct {
	ServiceMethod string
	Args          Serializer
	Reply         Serializer
	Error         error
	Done          chan *Call
}

func NewCall(serviceMethod string, args Serializer, reply Serializer) *Call {
	return &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          make(chan *Call, 1),
	}
}

func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	codec   ClientCodec
	sending sync.Mutex // guards the sending

	mu       sync.Mutex // protects following
	seq      uint64
	shutdown bool
	closing  bool
	pending  map[uint64]*Call
}

func NewClient(conn io.ReadWriteCloser) *Client {
	client := &Client{
		codec:   NewClientCodec(conn),
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client
}

func (c *Client) Call(serviceMethod string, args, reply Serializer) error {
	call := <-c.Go(serviceMethod, args, reply).Done
	return call.Error
}

func (c *Client) Go(serviceMethod string, args, reply Serializer) *Call {
	call := NewCall(serviceMethod, args, reply)
	c.send(call)
	return call
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.codec.Close()
}

func (c *Client) getSeq() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	seq := c.seq
	c.seq++
	return seq
}

func (c *Client) registerCall(seq uint64, call *Call) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pending[seq] = call
}

func (c *Client) getCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.pending[seq]
}

func (c *Client) removeCall(seq uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.pending, seq)
}

func (c *Client) send(call *Call) {
	c.sending.Lock()
	defer c.sending.Unlock()

	if c.shutdown || c.closing {
		call.Error = ErrShutdown
		call.done()
		return
	}

	req := &RequestHeader{
		ID:     c.getSeq(),
		Method: call.ServiceMethod,
	}
	c.registerCall(req.ID, call)

	body, err := call.Args.Marshal()
	if err != nil {
		call.Error = err
		c.removeCall(req.ID)
		return
	}
	req.Checksum = crc32.ChecksumIEEE(body)

	if err := c.writeRequest(req, body); err != nil {
		log.Println("rpc:failed to write request, err:", err)
		call.Error = err
		c.removeCall(c.seq)
		call.done()
	}
}

func (c *Client) receive() {
	var err error
	var response *ResponseHeader
	var data []byte
	for err == nil {
		response = new(ResponseHeader)
		data, err = c.readResponse(response)
		if err != nil {
			break
		}

		call := c.getCall(response.ID)
		c.removeCall(response.ID)
		if call != nil {
			if response.Error != "" {
				call.Error = errors.New(response.Error)
			} else if err := call.Reply.Unmarshal(data); err != nil {
				call.Error = err
			}
			call.done()
		}
	}
	// Terminate pending calls
	c.sending.Lock()
	c.mu.Lock()
	c.shutdown = true
	closing := c.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.mu.Unlock()
	c.sending.Unlock()
	if err != io.EOF && !closing {
		log.Println("rcc: client protocol error:", err)
	}
}

func (c *Client) writeRequest(req *RequestHeader, body []byte) error {
	if err := c.codec.WriteRequest(req, body); err != nil {
		return err
	}
	return nil
}

func (c *Client) readResponse(resp *ResponseHeader) ([]byte, error) {
	if err := c.codec.ReadResponseHeader(resp); err != nil {
		log.Println("rpc:failed to read response header, err:", err)
		return nil, err
	}

	data, err := c.codec.ReadResponseBody()
	if err != nil {
		log.Println("rpc:failed to read response body, err:", err)
		return nil, err
	}

	if resp.Checksum != crc32.ChecksumIEEE(data) {
		log.Println("rpc:response checksum mismatch")
		return nil, fmt.Errorf("rpc:response checksum mismatch")
	}

	return data, nil
}

type ClientCodec interface {
	WriteRequest(*RequestHeader, []byte) error
	ReadResponseHeader(*ResponseHeader) error
	ReadResponseBody() ([]byte, error)
	Close() error
}

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer
}

func NewClientCodec(conn io.ReadWriteCloser) ClientCodec {
	return &clientCodec{
		r: bufio.NewReader(conn),
		w: bufio.NewWriter(conn),
		c: conn,
	}
}

func (c *clientCodec) WriteRequest(req *RequestHeader, body []byte) error {
	if err := sendFrame(c.w, req.Marshal()); err != nil {
		log.Printf("rpc:failed to send request header, err is %s", err)
		return err
	}
	if err := sendFrame(c.w, body); err != nil {
		log.Printf("rpc:failed to send request body, err is %s", err)
		return err
	}

	c.w.(*bufio.Writer).Flush()
	return nil
}

func (c *clientCodec) ReadResponseHeader(r *ResponseHeader) error {
	data, err := recvFrame(c.r)
	if err != nil {
		log.Printf("rpc:failed to receive response header, err is %s", err)
		return err
	}

	return r.Unmarshal(data)
}

func (c *clientCodec) ReadResponseBody() (data []byte, err error) {
	data, err = recvFrame(c.r)
	if err != nil {
		log.Printf("rpc:failed to receive response body, err is %s", err)
	}
	return
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
