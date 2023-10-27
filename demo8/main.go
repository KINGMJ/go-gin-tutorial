package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	demo2()
}

type Person struct {
	Age      int       `form:"age" binding:"required,gt=10"`
	Name     string    `form:"name" binding:"omitempty"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func demo1() {
	r := gin.Default()
	r.GET("/validate", func(ctx *gin.Context) {
		var person Person
		if err := ctx.ShouldBind(&person); err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		ctx.String(http.StatusOK, fmt.Sprintf("%#v", person))
	})
	r.Run()
}

type Person2 struct {
	Age  int    `form:"age" binding:"required,gt=10"`
	Name string `form:"name" binding:"omitempty"`
	// 在参数 binding 上使用自定义的校验方法函数注册时候的名称
	BirthDate time.Time `form:"birthday" binding:"required,birth" time_format:"2006-01-02"`
}

func checkBirthDate(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return time.Now().After(value)
}

func demo2() {
	r := gin.Default()
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("birth", checkBirthDate)
	}
	// 接收请求
	r.GET("/validate", func(ctx *gin.Context) {
		var person Person2
		err := ctx.ShouldBindQuery(&person)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "success"})
	})
	r.Run()
}
