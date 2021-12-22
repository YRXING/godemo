package net

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

/*
Go标准包提供了对RPC的支持，但它只支持GO开发的服务器和客户端的交互，因为内部采用了Gob来编码。
编写规则：
1、函数必须是导出的 (首字母大写)
2、必须有两个导出类型的参数，
3、第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
4、函数还要有一个返回值 error
*/

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func HttpRPCServer() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP() //把服务注册到http协议上

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func HttpRPCClient(addr string) {
	client, err := rpc.DialHTTP("tcp", addr+":1234")
	if err != nil {
		fmt.Println(err)
	}

	//Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	fmt.Printf("result %d\n", reply)
}

func TcpRPCServer() {
	arith := new(Arith)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}

func TcpRPCClient(addr string) {
	client, err := rpc.Dial("tcp", addr+":1234")
	if err != nil {
		fmt.Println(err)
	}

	//Synchronous call
	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
}
