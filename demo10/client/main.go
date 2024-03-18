package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	done := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			body, err := getRequest()
			if err != nil {
				fmt.Println(err)
				close(done)
				return
			}
			fmt.Println(body)
		case <-done:
			return
		}
	}
}

func getRequest() (string, error) {
	// 并发访问接口
	res, err := http.Get("http://localhost:8080/ping")
	if err != nil || res.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("请求出错或者到达上限")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
