package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"jrpc_test/client"
	"net/http"
	"sync"
	"time"
)

type Work struct {
	sync.Mutex
	db *sql.DB
}

type (
	// SignPayload represent a request for vip-processor
	SignPayload struct {
		Name    string `json:"server_name"`
		Address string `json:"server_addr"`
		Port    int    `json:"server_port"`
	}
)

var W Work

//work可实现具体任务
func (w *Work) AddTask() error {
	fmt.Println("远程调用此方法")
	return nil
}

//在new 里初始化一些中间建如db等
func (w *Work) New() (*Work, error) {
	worker := &Work{}
	return worker, nil
}

type ControlService struct {
	Work *Work
}

func (c *ControlService) Ping(r *http.Request, in *client.PingMessage, out *client.PingMessage) error {
	out.Payload = in.Payload
	fmt.Println("我是远程服务我已启动，谢谢")
	return nil

}

func (c *Work) Sign(port int) {
	//获取心跳地址
	var addr string = "127.0.0.1:8888"
	go func(_addr string) {
		c := http.Client{
			Timeout: 250 * time.Millisecond,
		}
		for {
			_payload := &SignPayload{
				Name:    "sign",
				Address: "127.0.0.1",
				Port:    port,
			}

			_reqURL := "http://" + _addr + "/v1/registry/sign"
			blob, err := json.Marshal(_payload)
			if err != nil {
				return
			}
			req, err := http.NewRequest(http.MethodGet, _reqURL, bytes.NewBuffer(blob))
			if err != nil {
				fmt.Println("error", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Connection", "close")
			resp, err := c.Do(req)
			if err != nil {
				fmt.Println("心跳链接中断", err)
			}
			if resp != nil {
				if err := resp.Body.Close(); err != nil {
					panic("Close response body error")
				}
			}

			time.Sleep(time.Duration(3 * time.Second))
		}

	}(addr)
}
