去中心化的网络实现(1)TCP网络编程
============

> **作者:彭劲  版权所有,转载请注明来源**

# 为什么需要实现网络

经过前面的学习我们模拟的比特系统已经完成了区块链数据结构原型,工作量证明,UTXO交易,比特币地址和数字签名,并且优化了UTXO的交易机制。这些已经可以充分说明比特币核心技术大部分原理了,但前面完成的工作使得我们的比特币系统是"单节点"的,我们只有实现了网络才能更好地发挥前面工作实现的特性。这部分实现起来代码量相对较大,且与真正的比特币系统中的P2P网络实现是有差距的,但它有利于你对比特币网络技术的了解,并且为后面更好地学习以太坊网路奠定基础。
> 这部分的模拟实现代码需要使用到golang标准库中到net包实现网络,也就是tcp网络编程,所以我们需要先单独学习一下对应tcp网络编程的相关知识

# TCP服务器端的实现

使用golang标准库实现tcp网络编程,需要导入`net`包:

`net.Listen(network,port string)` 参数`network`传入`tcp`,表示当前协议为socket协议,并且监听指定的网络端口,参数`port`指定网络端口,如:`:3000`表示监听本地3000端口`localhost:3000`

```go
//如下程序等同于:listener,err := net.Listen("tcp",":3000")
listener,err := net.Listen("tcp","localhost:3000")
if err != nil {
	fmt.Println("can not connection localhost:3000")
	return //如果监听都无法实现则退出
}
```

这里值得一提的是`net.Listen`发生错误,程序就没有继续往下执行的必要性,因为这代表服务器根本无法建立,继续往下执行只会让后面的工作变得毫无意义

成功使用tcp监听网络的端口就代表我们创建了一个服务器,这时需要与客户端打通一个通信的数据链接通道,使用`listener.Accept()`不断等待客户端进行链接,如果没有检测到链接的发生,`listener.Accept()`会阻塞当前的代码继续往下执行

```go
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
		//响应链接,并产生通信链接,如果没有链接当前会阻塞
		conn,err := listener.Accept()
		defer listener.Close()
		if err != nil {
			fmt.Println("无法产生通信链接",err)
		}
		//处理当前通信链接
		go handlerConnection( conn )
	}
}
```

以上代码也许你会产生一个疑问,为什么`listener.Accept()`发生错误后程序继续往后执行?
因为服务器一旦创建成功户端可以接受多个客的链接,试想10个客户端链接服务器只有1个失败,我们是不可能终止程序的使得另外9个客户端退出的,所以这里发生错误只是某一个客户端的错误。

完成了上述的工作就需要考虑如何处理与客户端的通信了,答案是一个协程处理一个客户端通信链接。

```go
//处理当前通信链接
go handlerConnection( conn )
```

为什么要使用goroutine来实现协程处理客户端通信链接,目的是为了解决多客户端并发问题。1个协程服务1个客户端,多个客户端就无需等待`handlerConnection`函数完成执行。现在还需要看看`handlerConnection`的逻辑。

```go
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
```

这里最重要的是理解参数`conn`是保存了服务器端和客户端通信链接,这个链接如果不存在无法实现两者之间的通信。对于服务器端的解释就这么多[点击这里](https://github.com/pengjim520golang/blockchain-tutorial/blob/master/%E8%AF%BE%E7%A8%8B%E6%95%99%E6%9D%90/%E5%8E%BB%E4%B8%AD%E5%BF%83%E5%8C%96%E7%BD%91%E7%BB%9C%E7%9A%84%E5%AE%9E%E7%8E%B01/src/tcpServer/server.go)查看服务器端代码完整实现

# TCP客户端的实现

完成了服务器端的创建,客户端的实现相对简单,同样需要导入`net`包,使用`net.Dial`链接服务器端。

```go
package main 
import (
	"fmt"
	"net"
)
func main(){
	//链接服务器端
	conn,err := net.Dial("tcp","localhost:3000")
	if err != nil {
		fmt.Println("无法链接服务器")
		return 
	}
	//向服务器发送一条信息
	conn.Write([]byte("Hello BlockChainer.."))
}
```

`net.Dial(network,port string)` 参数`network`传入`tcp`与服务器端一致,参数`port`传入`localhost:3000`表示链接本地服务器监听的3000端口,`conn`是客户端和服务器端产生的通信链接通道,可以通过`conn.Write`向服务器端发送一条消息

> 注意：网络传输数据一般使用字节传送