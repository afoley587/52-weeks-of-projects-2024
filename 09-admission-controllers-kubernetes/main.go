package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ServeValidate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "validate",
	})
}
func ServeMutate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "mutate",
	})
}

func main() {
	log.Println("running...")
	r := gin.Default()
	r.GET("/validate", ServeValidate)
	r.GET("/mutate", ServeMutate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
