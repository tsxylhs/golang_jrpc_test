package main

import (
	"context"
	"flag"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"jrpc_test/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bind = flag.String("bind", ":9999", "server port")
var port = 9999

func main() {
	s := rpc.NewServer()
	//可传递参数配置文件参数等
	w, _ := service.W.New()
	s.RegisterCodec(json.NewCodec(), "application/json")
	err := s.RegisterService(&service.ControlService{Work: w}, "TEST")
	if err != nil {
		panic(err)
	}
	http.Handle("/rpc", s)

	srv := &http.Server{
		Addr: ":9999",
	}
	//保持心跳，
	service.W.Sign(port)
	go func() {
		//监听注册rpc服务
		if err = srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	//make channel 终止程序
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("程序终止")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务终止: %v", err)
	}

}
