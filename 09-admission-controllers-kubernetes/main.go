package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	admissionv1 "k8s.io/api/admission/v1"

	"github.com/afoley587/52-weeks-of-projects-2024/09-admission-controllers-kubernetes/pkg/admission"
)

func rawToAdmissionReview(c *gin.Context) (*admissionv1.AdmissionReview, error) {
	var req admissionv1.AdmissionReview

	if err := c.BindJSON(&req); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if req.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &req, nil

}

func ServeValidate(c *gin.Context) {

	req, err := rawToAdmissionReview(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Bad Request",
		})
		return
	}

	g := admission.GracefulAdmissionHandler{}
	review, err := g.AdmitValidation(req.Request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
		return
	}
	c.JSON(http.StatusOK, review)
}

func ServeMutate(c *gin.Context) {
	req, err := rawToAdmissionReview(c)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Bad Request",
		})
		return
	}

	g := admission.GracefulAdmissionHandler{}
	review, err := g.AdmitMutation(req.Request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
		return
	}
	c.JSON(http.StatusOK, review)
}

func main() {
	log.Println("running...")
	r := gin.Default()
	r.GET("/validate", ServeValidate)
	r.GET("/mutate", ServeMutate)
	r.Run() // listen and serve on 0.0.0.0:8080
}
