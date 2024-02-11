package main

import (
	"SmartHomeServer/app_client"
	"SmartHomeServer/smart_client"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go startItemServer()
	go startAppServer()
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

func startItemServer() {
	listen, err := net.Listen("tcp", "0.0.0.0:6502")
	fmt.Printf("物品服务端开启\n")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("物品建立连接失败, err:", err)
			continue
		}
		fmt.Println("物品建立了tcp连接,ip:", conn.RemoteAddr())
		// 对于每一个建立的tcp连接使用go关键字开启一个goroutine处理
		c := &smart_client.SmartClient{}
		go c.Start(conn)
	}
}

func startAppServer() {
	listen, err := net.Listen("tcp", "0.0.0.0:6503")
	fmt.Printf("APP服务端开启\n")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("APP建立连接失败, err:", err)
			continue
		}
		fmt.Println("APP建立了tcp连接,ip:", conn.RemoteAddr())
		// 对于每一个建立的tcp连接使用go关键字开启一个goroutine处理
		c := &app_client.AppClient{}
		go c.Start(conn)
	}
}
