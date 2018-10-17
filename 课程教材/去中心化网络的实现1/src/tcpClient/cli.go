package main 

import (
	"fmt"
	"net"
)

func main(){
	conn,err := net.Dial("tcp","localhost:3000")
	if err != nil {
		fmt.Println("无法链接服务器")
		return 
	}
	conn.Write([]byte("Hello BlockChainer.."))
}
