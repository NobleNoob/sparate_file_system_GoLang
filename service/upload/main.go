package main

import (
	"fmt"
	"net/http"
	cfg "filestore-server/config"
)

func main() {

	// 监听端口
	fmt.Println("上传服务正在启动, 监听端口:%s...",cfg.UploadServiceHost)
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}