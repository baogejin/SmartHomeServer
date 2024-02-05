package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go startServer()
	exit := make(chan bool, 1)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-ch
		fmt.Println("server exit")

		exit <- true
	}()
	<-exit
}

func startServer() {
	listen, err := net.Listen("tcp", "0.0.0.0:6502")
	fmt.Printf("服务端: %T=====\n", listen)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		fmt.Println("当前建立了tcp连接")
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		// 对于每一个建立的tcp连接使用go关键字开启一个goroutine处理
		go process(conn)
	}
}

func process(conn net.Conn) {
	// 函数执行完之后关闭连接
	defer conn.Close()
	// 输出主函数传递的conn可以发现属于*TCPConn类型, *TCPConn类型那么就可以调用*TCPConn相关类型的方法, 其中可以调用read()方法读取tcp连接中的数据
	fmt.Printf("服务端: %T\n", conn)
	conn.Write([]byte("hello\n"))
	for {
		var buf [128]byte
		// 将tcp连接读取到的数据读取到byte数组中, 返回读取到的byte的数目
		n, err := conn.Read(buf[:])
		if err != nil {
			// 从客户端读取数据的过程中发生错误
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("服务端收到客户端发来的数据：", recvStr)
		// 向当前建立的tcp连接发送数据, 客户端就可以收到服务端发送的数据
		conn.Write([]byte(recvStr))
	}
}
