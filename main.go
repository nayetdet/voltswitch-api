package main

import (
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNoContent, nil)
	})

	server.POST("/shutdown", func(ctx *gin.Context) {
		command := strings.TrimSpace(os.Getenv("SHUTDOWN_COMMAND"))
		if command == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "SHUTDOWN_COMMAND environment variable is not set",
			})

			return
		}

		commandParts := strings.Fields(command)
		err := exec.Command(commandParts[0], commandParts[1:]...).Run()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	})

	if err := server.Run(":8000"); err != nil {
		panic(err)
	}
}
