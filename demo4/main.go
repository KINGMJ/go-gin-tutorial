package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LoginForm struct {
	User     string `form:"user" xml:"user" json:"user"`
	Password string `form:"password" xml:"password" json:"password"`
}

func main() {
	demo3()
}

func demo1() {
	r := gin.Default()
	r.GET("/login", func(ctx *gin.Context) {
		var form LoginForm
		if ctx.ShouldBind(&form) == nil {
			ctx.JSON(http.StatusOK, form)
		}
	})
	r.Run()
}

func demo2() {
	r := gin.Default()
	r.POST("/login", func(ctx *gin.Context) {
		var form, form1, form2, form3 LoginForm
		if ctx.ShouldBindWith(&form, binding.Form) == nil {
			ctx.JSON(http.StatusOK, fmt.Sprintf("binding.Form %s", form))
		}
		if ctx.ShouldBindWith(&form1, binding.Query) == nil {
			ctx.JSON(http.StatusOK, fmt.Sprintf("binding.Query %s", form1))
		}
		if ctx.ShouldBindWith(&form2, binding.JSON) == nil {
			ctx.JSON(http.StatusOK, fmt.Sprintf("binding.JSON %s", form2))
		}
		if ctx.ShouldBindWith(&form3, binding.XML) == nil {
			ctx.JSON(http.StatusOK, fmt.Sprintf("binding.XML %s", form3))
		}
	})
	r.Run()
}

type Person struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func demo3() {
	r := gin.Default()
	r.GET("/:name/:id", func(ctx *gin.Context) {
		var person Person
		if err := ctx.ShouldBindUri(&person); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"name": person.Name, "uuid": person.ID})
	})
	r.Run()
}
