package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		//等待客户端conn发送信息，如果客户端没有发送，那么协程就阻塞在这里
		fmt.Printf("服务器在等待客户端%s 发送信息\n", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("客户端退出 err=%v", err)
			return //!!!
		}
		//显示客户端发送的内容到服务器的终端
		fmt.Print(string(buf[:n]))
	}

}
func main() {
	fmt.Println("服务器开始监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	defer listen.Close() //延时关闭listen
	//循环等待客户端来连接我
	for {
		fmt.Println("等待客户端来连接...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept() err=", err)
		} else {
			fmt.Printf("Accept() success con=%v 客户端ip=%v\n", conn, conn.RemoteAddr().String())
		}
		//这里准备一个协程为客户端服务
		go process(conn)
	}
}
