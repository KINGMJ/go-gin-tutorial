package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	demo2()
}

// 单文件上传
func demo1() {
	r := gin.Default()
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MB)
	r.MaxMultipartMemory = 8 << 20

	r.POST("/upload", func(ctx *gin.Context) {
		// 单文件
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		log.Println(file.Filename)
		// 上传文件至指定的完整文件路径，默认为当前文件夹
		ctx.SaveUploadedFile(file, file.Filename)
		ctx.JSON(http.StatusOK, file)
	})
	r.Run()
}

// 多文件上传
func demo2() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		files := form.File["files"]
		for _, file := range files {
			if err := ctx.SaveUploadedFile(file, file.Filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}
		}
		ctx.JSON(http.StatusOK, files)
	})
	r.Run()
}
