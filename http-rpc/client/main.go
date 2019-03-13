package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	// 提供 -s 选项指定 RPC Server 的监听地址和端口
	serverAddress := flag.String("s", "127.0.0.1:1234", "Assign the address and port of RPC Server.")
	flag.Parse()

	// 连接 RPC 服务端
	client, err := rpc.DialHTTP("tcp", *serverAddress)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// 同步调用
	args := &Args{3, 4}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	// 异步调用
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	// 在上面的异步调用命令发出之后，客户端可以去执行其它的任务
	replyCall := <-divCall.Done // divCall.Done 是一个 channel，如果 RPC 异步调用有结果了，就可以从 chan 中读取到值
	// check errors, print, etc.
	if err := replyCall.Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quotient.Quo, quotient.Rem)
}
