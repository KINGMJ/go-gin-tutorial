package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	demo7()
}

// 基本路由
func demo1() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.Run()
}

// 路由参数
func demo2() {
	r := gin.Default()
	r.GET("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "Hello %s", name)
	})
	r.GET("/user/:name/*action", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, ctx.Params)
	})
	r.GET("/user/settings", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "settings ok")
	})
	r.Run()
}

// 路由分组
func demo3() {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/login", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "v1 login endpoint")
		})
		v1.GET("/submit", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "v1 submit endpoint")
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/login", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "v2 login endpoint")
		})
		v2.GET("/submit", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "v2 submit endpoint")
		})
	}

	r.GET("/login", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "login")
	})

	r.Run()
}

// url query 参数
func demo4() {
	r := gin.Default()
	r.GET("/user", func(ctx *gin.Context) {
		name := ctx.Query("name")
		page := ctx.DefaultQuery("page", "0")
		id, ok := ctx.GetQuery("id")
		if !ok {
			fmt.Println("参数不存在")
		}
		fmt.Println(id)
		fmt.Println(name)
		fmt.Println(page)
	})

	r.Run()
}

// post form 参数
func demo5() {
	r := gin.Default()
	r.HEAD("/user", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		age := ctx.DefaultPostForm("age", "12")
		sex, ok := ctx.GetPostForm("sex")
		if !ok {
			fmt.Println("参数不存在")
		}
		fmt.Println(name)
		fmt.Println(age)
		fmt.Println(sex)
	})
	r.Run()
}

// 数组和对象参数
func demo6() {
	r := gin.Default()
	r.POST("/user", func(ctx *gin.Context) {
		// names := ctx.PostFormArray("name")
		names := ctx.PostFormMap("name")
		ctx.JSON(http.StatusOK, names)
	})
	r.Run()
}

// 获取客户端ip
func demo7() {
	r := gin.Default()
	r.POST("/ip", func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		fmt.Println(ip)
	})
	r.Run()
}
