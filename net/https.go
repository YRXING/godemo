package net

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func HttpsClient() {
	conf := &tls.Config{
		//不会校验证书以及证书中的主机名是否和服务器一致
		InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", "127.0.0.1:443", conf)

	defer conn.Close()

	n, _ := conn.Write([]byte("hello\n"))
	buf := make([]byte, 100)
	n, _ = conn.Read(buf)
	fmt.Println(string(buf[:n]))
}

func HttpsServer() {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
	}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(msg)
		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
