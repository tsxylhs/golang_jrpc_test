package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jrpc_test/client"
	"log"
)

func main() {
	router := gin.Default()
	//router := gin.New()

	router.GET("/v1/registry/sign", func(context *gin.Context) {
		log.Println(">>>> hello jrpc <<<<")

		_decoderClient := client.New("http://127.0.0.1:9999/rpc")
		resq, err := _decoderClient.Ping("我的接口通了")
		if err != nil {
			fmt.Println("接口调用失败")
		} else {
			fmt.Println("远程调用成功", resq)
		}
		context.JSON(200, gin.H{
			"code":    200,
			"success": true,
		})
	})
	// 指定地址和端口号
	router.Run("localhost:8888")
}
