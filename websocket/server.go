package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} //user default options

func EchoServer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		msgtype, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("server recv：%s", msg)
		err = conn.WriteMessage(msgtype, msg)
		if err != nil {
			log.Println("write：", err)
			break
		}
	}
}

//server端是一个http 服务器，监听8080端口。当接收到连接请求后，将连接使用的http协议升级为websocket协议。
//后续通信过程中，使用websocket进行通信
func main() {
	flag.Parse()
	http.HandleFunc("/echo", EchoServer)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
