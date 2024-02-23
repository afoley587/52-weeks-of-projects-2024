package routes

import (
	"net/http"

	"github.com/afoley587/52-weeks-of-projects/07-golang-gin-langchain/chains"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type GenerateVacationIdeaRequest struct {
	FavoriteSeason string   `json:"favorite_season"`
	Hobbies        []string `json:"hobbies"`
	Budget         int      `json:"budget"`
}

type GenerateVacationIdeaResponse struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
}

type GetVacationIdeaResponse struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}

func generateVacation(r GenerateVacationIdeaRequest) GenerateVacationIdeaResponse {
	// call the chain here
	return GenerateVacationIdeaResponse{Id: uuid.New(), Completed: false}
}

func getVacation(id uuid.UUID) GetVacationIdeaResponse {
	// search the chains here
	v, _ := chains.GetVacationFromDb(id)
	return GetVacationIdeaResponse{Id: v.Id, Completed: v.Completed, Idea: v.Idea}
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
		id, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			resp := getVacation(id)
			c.JSON(http.StatusOK, resp)

		}
	})
	return router
}
