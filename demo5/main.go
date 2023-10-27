package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	demo6()
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

	r.GET("/xml", func(ctx *gin.Context) {
		ctx.XML(http.StatusOK, gin.H{"message": "hey", "status": 200})
	})

	r.GET("/yaml", func(ctx *gin.Context) {
		ctx.YAML(http.StatusOK, gin.H{"message": "hey", "status": 200})
	})
	r.Run()
}

func demo2() {
	r := gin.Default()
	r.GET("/file", func(ctx *gin.Context) {
		ctx.File("/Users/jerrymei/Desktop/_DSC3118.jpg")
	})
	r.GET("/download", func(ctx *gin.Context) {
		ctx.FileAttachment("/Users/jerrymei/Desktop/_DSC3118.jpg", "1.jpg")
	})
	r.Run()
}

func demo3() {
	r := gin.Default()
	r.GET("/file", func(ctx *gin.Context) {
		ctx.File("/Users/jerrymei/Desktop/_DSC3118.jpg")
	})
	r.GET("/download", func(ctx *gin.Context) {
		ctx.FileAttachment("/Users/jerrymei/Desktop/_DSC3118.jpg", "1.jpg")
	})
	r.Run()
}

func demo4() {
	r := gin.Default()
	r.Static("/assets", "/Users/jerrymei/Desktop/img")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFS("/static", http.Dir("/Users/jerrymei/Desktop/img"))
	r.Run()
}

func demo5() {
	r := gin.Default()
	r.GET("/test", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
	})

	r.POST("/test", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/foo")
	})

	r.GET("/foo", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	r.GET("/test2", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/test3"
		r.HandleContext(ctx)
	})

	r.GET("/test3", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	r.Run()
}

func demo6() {
	r := gin.Default()
	// 异步
	r.GET("/async", func(ctx *gin.Context) {
		copyContext := ctx.Copy()
		go func() {
			time.Sleep(3 * time.Second)
			log.Println("异步执行：" + copyContext.Request.URL.Path)
		}()
	})

	// 同步
	r.GET("/async1", func(ctx *gin.Context) {
		time.Sleep(3 * time.Second)
		log.Println("同步执行：" + ctx.Request.URL.Path)
	})

	r.Run()
}
