package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	demo4()
}

func demo1() {
	r := gin.New()
	// 通过use设置全局中间件
	// 设置日志中间件，主要用于打印请求日志
	r.Use(gin.Logger())

	// 设置Recovery中间件，主要用于拦截panic错误，不至于导致进程崩掉
	r.Use(gin.Recovery())

	// 忽略后面代码
	r.GET("middleware", func(ctx *gin.Context) {
		ctx.String(200, "hello world")
	})

	r.Run()
}

// 自定义中间件
func demo2() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/log", func(ctx *gin.Context) {
		example := ctx.MustGet("example").(string)
		ctx.String(200, example)
		log.Println(example)
	})

	r.Run()
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		// 可以通过上下文对象，设置一些依附在上下文对象里面的键/值数据
		ctx.Set("example", "123456")

		// 在这里处理请求到达控制器函数之前的逻辑

		// 调用下一个中间件，或者控制器处理函数，具体得看注册了多少个中间件。
		ctx.Next()

		// 在这里可以处理请求返回给用户之前的逻辑
		latency := time.Since(t)
		log.Printf("延迟：%s\n", latency)

		// 例如，查询请求状态吗
		status := ctx.Writer.Status()
		log.Printf("状态码：%d\n", status)

	}
}

// 默认启用Recovery中间件，panic不会崩溃
func demo3() {
	r := gin.Default()
	r.GET("/recovery", func(ctx *gin.Context) {
		panic("errors")
	})
	r.Run()
}

func demo4() {
	r := gin.Default()
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "123",
		"foo":   "bar",
	}))

	authorized.GET("/login", func(ctx *gin.Context) {
		user := ctx.MustGet(gin.AuthUserKey).(string)
		ctx.JSON(200, user)
	})

	r.Run()
}
