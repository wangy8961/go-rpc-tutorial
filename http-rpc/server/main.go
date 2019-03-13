package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
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
	rpc.HandleHTTP()    // RPC 会处理客户端发来的 HTTP 请求

	// 启动 HTTP 服务器，监听在 1234 端口
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)

	/* 或者:
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	*/
}
