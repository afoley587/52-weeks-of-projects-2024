package main

import (
	"github.com/afoley587/52-weeks-of-projects/07-golang-gin-langchain/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.GetVacationRouter(r)
	r.Run()
}
