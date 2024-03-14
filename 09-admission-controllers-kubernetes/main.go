package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/afoley587/52-weeks-of-projects/09-admission-controllers-kubernetes/admission"
)

func ServeValidate(c *gin.Context) {
	g := admission.GracefulAdmissionHandler{}
	review, err := g.Validate()
	// c.JSON(200, gin.H{
	// 	"message": "validate",
	// })
}
func ServeMutate(c *gin.Context) {
	g := admission.GracefulAdmissionHandler{}
	review, err := g.Mutate()
	// c.JSON(200, gin.H{
	// 	"message": "mutate",
	// })
}

func main() {
	log.Println("running...")
	r := gin.Default()
	r.GET("/validate", ServeValidate)
	r.GET("/mutate", ServeMutate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
