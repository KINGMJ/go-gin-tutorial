package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	demo1()
}

func demo1() {
	r := gin.Default()
	r.GET("/string", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
		ctx.String(http.StatusOK, "欢迎访问%s，你是%s", "www.baidu.com", "jack")
	})

	r.GET("/json", func(ctx *gin.Context) {
		// 使用一个结构体
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 在 JSON 中变成了 "user"
		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
		ctx.JSON(http.StatusOK, msg)
		// gin.H 是 map[string]interface{} 的一种快捷方式
		ctx.JSON(http.StatusOK, gin.H{"message": "hey", "status": 200})
	})
	r.Run()
}
