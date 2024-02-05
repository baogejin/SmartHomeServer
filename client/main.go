package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 连接到服务端建立的tcp连接
	conn, err := net.Dial("tcp", "192.168.31.136:6502")
	// 输出当前建Dial函数的返回值类型, 属于*net.TCPConn类型
	fmt.Printf("客户端: %T\n", conn)
	if err != nil {
		// 连接的时候出现错误
		fmt.Println("err :", err)
		return
	}
	// 当函数返回的时候关闭连接
	defer conn.Close()
	// 获取一个标准输入的*Reader结构体指针类型的变量
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 调用*Reader结构体指针类型的读取方法
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		// 去除掉\r \n符号
		inputInfo := strings.Trim(input, "\r\n")
		// 判断输入的是否是Q, 如果是Q则退出
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			return
		}
		_, err = conn.Write([]byte(inputInfo)) // 发送数据
		if err != nil {
			return
		}
		buf := [512]byte{}
		// 读取服务端发送的数据
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println("客户端接收服务端发送的数据: ", string(buf[:n]))
	}
}
