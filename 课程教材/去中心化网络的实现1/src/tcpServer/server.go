package main

import (
	"fmt"
	"net"
)
func main(){
	listener,err := net.Listen("tcp","localhost:3000")
	if err != nil {
		fmt.Println("can not connection localhost:3000")
		return
	}
	//不断等待客户端的链接
	for {
		//响应链接,并产生通信链接
		conn,err := listener.Accept()
		defer listener.Close()
		if err != nil {
			fmt.Println("无法产生通信链接",err)
		}
		//处理当前通信链接
		go handlerConnection( conn )
	}
}

func handlerConnection( conn net.Conn ){
	defer conn.Close()
	//不断等待客户端的发送
	for {
		//读取客户端发送过来数据的容器,表示最多读取1024字节
		buf := make([]byte,1024)
		n,err := conn.Read( buf )
		if err != nil {
			fmt.Println("客户端退出...")
			return 
		}
		message := string(buf[:n])
		fmt.Printf("%s:%s\n",conn.RemoteAddr().String(), message)
	}
}
