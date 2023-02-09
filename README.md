# drpc

> 本仓库是drpc的go语言的实现。drpc需与 [dgen](https://github.com/fengluodb/dgen) 配合使用。

## 使用方法

首先，新建一个idl文件，然后根据[IDL语法](https://github.com/fengluodb/dgen#idl-%E8%AF%AD%E6%B3%95)定义服务。

这里定义一个简单的服务：

```protobuf
message HelloRequest {
    seq=1 string name;
}

message HelloReply {
    optional seq=1 string reply;
}

service HelloWorld {
    SayHello(HelloRequest) return (HelloReply);
}
```

使用dgen工具，如`dgen -f helloworld.idl -l go`, 可以生成一个helloworld目录。目录下的hellowrold.go中定义了 `HelloRequest`、`HelloReply` 和它们的序列化、反序列化方法，hellowrold.dprc.go中定义了包括服务注册在内的一些桩代码。用户只需要编写少量的服务端代码，客户端就可以实现调用。

**server**：
```go
type helloService struct {
}

func (h *helloService) SayHello(args *dgen.HelloRequest, reply *dgen.HelloReply) error {
	reply.Reply = fmt.Sprintf("Hello %s, welcome to drpc", args.Name)
	return nil
}

func main() {
	server := drpc.NewServer()
	dgen.RegisterHelloWorldService(server, "Hello", new(helloService)) // 此函数由编译器生成

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}
	server.Serve(listener)
}
```

**client**:
```go
func main() {
	client, err := drpc.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("err:", err)
	}

	args := &dgen.HelloRequest{
		Name: "fengluodb",
	}
	reply := &dgen.HelloReply{}

	if err := client.Call("Hello.SayHello", args, reply); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("reply:", reply.Reply)
}
```
在本仓库的 [hellowrold](https://github.com/fengluodb/drpc/tree/main/example/helloworld) 目录下有该示例，其中`defalut`和`json`代表不同的序列化方式。

## 传输协议

**请求头**:
```go
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
```
`ID`为每个请求的唯一标识，`Method`为调用的方法名，`Checksum`用于检查request body传输过程中是否发生错误。


**响应头**
```go
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
```
`ID`为每个请求的唯一标识，`Error`代表函数调用时是否发生错误（如果`Error`为空，代表没有错误），`Checksum`用于检查response body传输过程中是否发生错误。