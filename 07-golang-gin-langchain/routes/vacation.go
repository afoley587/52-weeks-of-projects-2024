package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenerateVacationIdeaRequest struct {
	FavoriteSeason string   `json:"favorite_season"`
	Hobbies        []string `json:"hobbies"`
	Budget         int      `json:"budget"`
}

type GenerateVacationIdeaResponse struct {
	Id        int  `json:"id"`
	Completed bool `json:"completed"`
}

type GetVacationIdeaResponse struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Idea      string `json:"idea"`
}

func generateVacation(r GenerateVacationIdeaRequest) GenerateVacationIdeaResponse {
	return GenerateVacationIdeaResponse{Id: 1, Completed: false}
}

func GetVacationRouter(router *gin.Engine) *gin.Engine {

	registrationRoutes := router.Group("/vacation")
	registrationRoutes.POST("/create", func(c *gin.Context) {
		var req GenerateVacationIdeaRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			c.JSON(http.StatusOK, generateVacation(req))
		}
	})
	registrationRoutes.GET("/:id", func(c *gin.Context) {
		sid := c.Param("id")
		id, _ := strconv.Atoi(sid)
		resp := GetVacationIdeaResponse{Id: id, Idea: "", Completed: true}
		c.JSON(http.StatusOK, resp)

		/*
			var req GenerateVacationIdeaRequest
			err := c.BindJSON(&req)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Bad Request",
				})
			} else {
				c.JSON(http.StatusOK, generateVacation(req))
			}
		*/
	})
	return router
}
