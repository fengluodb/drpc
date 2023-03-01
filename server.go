package drpc

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var _ ServerCodec = (*serverCodec)(nil)

type Handler func(args []byte) ([]byte, error)

type service struct {
	methodMap map[string]Handler
}

func NewService() *service {
	return &service{
		methodMap: make(map[string]Handler),
	}
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go s.ServeConn(conn)
	}
}

func (s *Server) ServeConn(conn net.Conn) {
	codec := NewServerCodec(conn)
	s.ServeCodec(codec)
}

func (s *Server) ServeCodec(codec ServerCodec) {
	defer codec.Close()
	for {
		req, handler, args, err := s.readRequest(codec)
		if err != nil {
			log.Println("rpc:failed to read request, err:", err)
			break
		}

		if err := s.call(req, codec, handler, args); err != nil {
			log.Println("rpc:failed to send response, err:", err)
			break
		}
	}
}

func (s *Server) readRequest(codec ServerCodec) (req *RequestHeader, handler Handler, args []byte, err error) {
	req = new(RequestHeader)
	err = codec.ReadRequestHeader(req)
	if err != nil {
		return
	}

	dot := strings.LastIndex(req.Method, ".")
	if dot < 0 {
		err = fmt.Errorf("rpc: service/method request ill-formed: %s", req.Method)
		return
	}
	serviceName := req.Method[:dot]
	methodName := req.Method[dot+1:]

	// look for the method
	svci, ok := s.serviceMap.Load(serviceName)
	if !ok {
		err = fmt.Errorf("can't find service:%s", serviceName)
		return
	}
	svc := svci.(*service)
	handler, ok = svc.methodMap[methodName]
	if !ok {
		err = fmt.Errorf("can't find method:%s", methodName)
		return
	}

	args, err = codec.ReadRequestBody()
	if req.Checksum != crc32.ChecksumIEEE(args) {
		err = fmt.Errorf("request checksum mismatch")
	}
	return
}

func (s *Server) call(req *RequestHeader, codec ServerCodec, handler Handler, args []byte) error {
	resp := new(ResponseHeader)
	reply, err := handler(args)

	resp.ID = req.ID
	if err != nil {
		resp.Error = err.Error()
	}
	resp.Checksum = crc32.ChecksumIEEE(reply)
	if err := codec.WriteResponse(resp, reply); err != nil {
		log.Printf("rpc:failed to send response, err:%s", err)
		return err
	}
	return nil
}

func RegisterService(s *Server, serviceMethodName string, method Handler) error {
	dot := strings.LastIndex(serviceMethodName, ".")
	if dot == -1 {
		log.Println("serviceMethod bust be the format of serviceName.methodName")
		return fmt.Errorf("serviceMethod bust be the format of serviceName.methodName")
	}
	serviceName := serviceMethodName[:dot]
	methodName := serviceMethodName[dot+1:]

	svci, ok := s.serviceMap.Load(serviceName)
	if !ok {
		svci = NewService()
		s.serviceMap.Store(serviceName, svci)
	}
	svc := svci.(*service)

	if _, ok := svc.methodMap[methodName]; ok {
		log.Printf("rpc:%s has been registered", serviceMethodName)
		return fmt.Errorf("%s has been registered", serviceMethodName)
	}
	svc.methodMap[methodName] = method
	log.Printf("rpc:register %s successfully", serviceMethodName)

	return nil
}

type ServerCodec interface {
	ReadRequestHeader(*RequestHeader) error
	ReadRequestBody() ([]byte, error)
	WriteResponse(*ResponseHeader, []byte) error
	Close() error
}

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	closed bool
}

func NewServerCodec(conn io.ReadWriteCloser) ServerCodec {
	return &serverCodec{
		r: bufio.NewReader(conn),
		w: bufio.NewWriter(conn),
		c: conn,
	}
}

func (s *serverCodec) ReadRequestHeader(r *RequestHeader) error {
	data, err := recvFrame(s.r)
	if err != nil {
		log.Printf("rpc:failed to receive request header, err is %s", err)
		return err
	}

	return r.Unmarshal(data)
}

func (s *serverCodec) ReadRequestBody() ([]byte, error) {
	return recvFrame(s.r)
}

func (s *serverCodec) WriteResponse(resp *ResponseHeader, body []byte) error {
	if err := sendFrame(s.w, resp.Marshal()); err != nil {
		log.Printf("rpc:failed to send request header, err is %s", err)
		return err
	}
	if err := sendFrame(s.w, body); err != nil {
		log.Printf("rpc:failed to send request body, err is %s", err)
		return err
	}

	return s.w.(*bufio.Writer).Flush()
}

func (s *serverCodec) Close() error {
	if s.closed {
		return nil
	}

	s.closed = true
	return s.c.Close()
}
