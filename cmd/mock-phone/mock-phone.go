package main

import (
	passwordentry "Lunarisnia/strongest-pin-finder/internal/password_entry"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Oke",
		})
	})

	r.POST("/login", passwordentry.ProcessPasswordEntry)
	r.Run(":3000")
}
