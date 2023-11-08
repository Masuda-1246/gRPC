
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Masuda-1246/gRPC/api/todopb"
	"google.golang.org/grpc"
)

func checkErr(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type server struct{}

func main() {
	fmt.Println("=======Todo API Start====")

	// 50051ポートで起動する(docker-compose.ymlと合わせる)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	checkErr("Failed to listen: %v", err)

	// grpcのサーバーを作成
	s := grpc.NewServer()
	// TodoServiceServerにマウント
	todopb.RegisterTodoServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		err = s.Serve(lis)
		checkErr("faild to server: %v", err)
	}()

	// Ctrl + cでプログラムから抜けられるようにする
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// シグナルが受け取れるまで、ブロック
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("=====Todo API End====")
}