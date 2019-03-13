package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

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

func main() {
	arith := new(Arith)
	rpc.Register(arith) // 将 Arith 类型的变量 arith 的所有方法注册到 RPC 中，比如 Multiply() 和 Divide()

	// 创建监听套接字
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	for {
		// 为每个客户端连接创建连接套接字
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Println("connection error: ", err)
			continue
		}
		// 启动 JSON-RPC Server，处理客户端发来的 TCP 请求
		jsonrpc.ServeConn(conn)
	}
}
