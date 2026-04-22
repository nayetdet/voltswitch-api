package main

import (
	"net/http"
	"os/exec"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNoContent, nil)
	})

	server.POST("/shutdown", func(ctx *gin.Context) {
		output, err := exec.Command("shutdown", "-h", "now").CombinedOutput()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"output": string(output),
				"error":  err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	})

	if err := server.Run(":8000"); err != nil {
		panic(err)
	}
}
