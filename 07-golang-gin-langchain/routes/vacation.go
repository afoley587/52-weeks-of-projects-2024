package routes

import (
	"net/http"

	"github.com/afoley587/52-weeks-of-projects/07-golang-gin-langchain/chains"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

/*

Now that our request and response schemas are out
of the way, we can talk about the router.

The `GetVacationRouter` function takes a gin router
as input and adds a new router grouping to it with a
path prefix of `/vacation`. So any endpoints we add to
the router will get that `/vacation` prefix. We then
add two endpoints:

1. POST `/create` to create a new vacation idea for us
2. GET `/:id` which gets an idea by ID

From a high level, the `/create` endpoint will kick
off a goroutine which invokes langchain and openAI. It
will return a `GenerateVacationIdeaResponse` to the caller
so that they can check back on the status of it later. They
check the status of that idea with the `/:id` endpoint. This
will return a `GetVacationIdeaResponse`. If the idea is finished
generating, it will contain an id, an idea, and the completed flag
will be set to true. Otherwise, it will contain an id, an empty idea,
and the completed flag will be set to false.

So, in more depth, what happens when we send a POST to `/create`? First,
it will validate the request, i.e. making sure that the request that
was sent in is valid. If it is, it will call `generateVacation` with
the deserialized request. `generateVacation` will generate a new UUID
for the idea and invoke the langchain chain in a goroutine with
`go chains.GeneateVacationIdeaChange(id, r.Budget, r.FavoriteSeason, r.Hobbies)`.
We will talk about what this chain is doing in the next section but, for now,
let's just say that it goes and generates a vacation idea with the given parameters.
After it starts that chain, it returns a `GenerateVacationIdeaResponse` to the
caller with the ID field set. We should also note that it's important to us to
put this on a goroutine because we want our responses to be snappy. Langchain
might take a few seconds to generate an actual idea so we don't want clients
to register things like timeouts when waiting for a response.

How about the GET to `/:id`? Well, this one's a bit simpler. It also first
validates the request by making sure that the ID is a valid UUID. It then goes
and tries to query the in-memory DB (a map) that we keep for vacation idea
book-keeping. If the ID doesn't exist, we will return a 404 to the caller.
If it does, we will turn the `Vacation` object into a `GetVacationIdeaResponse`
containing all the releveant data (ID, idea, and whether it's finished or not).
*/

func generateVacation(r GenerateVacationIdeaRequest) GenerateVacationIdeaResponse {
	// First, generate a new UUID for the idea
	id := uuid.New()

	// Then invoke the GeneateVacationIdeaChange method of the chains package
	// passing through all of the parameters from the user
	go chains.GeneateVacationIdeaChange(id, r.Budget, r.FavoriteSeason, r.Hobbies)
	return GenerateVacationIdeaResponse{Id: id, Completed: false}
}

func getVacation(id uuid.UUID) (GetVacationIdeaResponse, error) {
	// Search the chains database for the ID requested by the user
	v, err := chains.GetVacationFromDb(id)

	// If the ID didn't exist, handle the error
	if err != nil {
		return GetVacationIdeaResponse{}, err
	}

	// Otherwise, return the vacation idea to the caller
	return GetVacationIdeaResponse{Id: v.Id, Completed: v.Completed, Idea: v.Idea}, nil
}

func GetVacationRouter(router *gin.Engine) *gin.Engine {

	// Add a new router group to the gin router
	registrationRoutes := router.Group("/vacation")

	// Handle the POST to /create
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

	// Handle the GET to /:id
	registrationRoutes.GET("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			resp, err := getVacation(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Id Not Found",
				})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		}
	})

	// Return the updated router
	return router
}
