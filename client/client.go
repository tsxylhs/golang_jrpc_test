package client

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	"net/http"
	"time"
)

//声明clent 链接客户端地址
type Client struct {
	Address string
}

//将client 地址赋值
func New(addr string) *Client {
	return &Client{
		Address: addr,
	}
}

//jrp实现
func (c *Client) jrpc(method string, in interface{}, out interface{}) error {
	message, err := json.EncodeClientRequest(method, in)
	if err != nil {
		return err
	}
	//封装Http请求作物rpc 的载体
	fmt.Println("c.address", c.Address)
	req, err := http.NewRequest(http.MethodPost, c.Address, bytes.NewBuffer(message))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求执行失败")
		return err
	}
	defer resp.Body.Close()
	return json.DecodeClientResponse(resp.Body, out)
}

type PingMessage struct {
	Payload string
}

//测试远程服务是否启动
func (c *Client) Ping(message string) (string, error) {
	in := PingMessage{Payload: message}
	var out PingMessage
	err := c.jrpc("TEST.Ping", in, &out)
	if err != nil {
		fmt.Println("接口调用失败", err)
		return "", err
	}
	return out.Payload, nil
}

//其他方法加入开箱加入即可
