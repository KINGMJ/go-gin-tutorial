package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	demo2()
}

func demo1() {
	r := gin.Default()

	r.GET("/cookie", func(ctx *gin.Context) {
		ctx.SetCookie("site_cookie", "cookievalue", 3600, "/", "localhost", false, true)
	})

	r.GET("/getcookie", func(ctx *gin.Context) {
		data, err := ctx.Cookie("site_cookie")
		if err != nil {
			ctx.String(http.StatusOK, "not found!")
			return
		}
		ctx.String(http.StatusOK, data)
	})

	r.GET("/deletecookie", func(ctx *gin.Context) {
		ctx.SetCookie("site_cookie", "cookievalue", -1, "/", "localhost", false, true)
	})
	r.Run()
}

// cookie 练习 模拟实现权限验证中间件。
// 有2个路由，login和home。 login用于设置cookie；home是访问查看信息的请求
// 在请求home之前，先跑中间件代码，检验是否存在cookie
func demo2() {
	r := gin.Default()
	r.GET("/login", func(ctx *gin.Context) {
		ctx.SetCookie("abc", "123", 60, "/", "localhost", false, true)
		ctx.String(200, "Logic success!")
	})

	r.GET("/home", AuthMiddleWare(), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"data": "home"})
	})

	r.Run()
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("abc")
		if err != nil || cookie != "123" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
