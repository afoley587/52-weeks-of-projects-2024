package main

/*
Running

go run main.go

Example Curls:

curl -X POST -H"Content-type: application/json" -d'{"favorite_season": "summer", "hobbies": ["surfing","running"], "budget":1000}' http://localhost:8080/vacation/create
curl -X GET -H"Content-type: application/json" http://localhost:8080/vacation/ed738008-e56c-450e-a3ff-83bf8e85d167
*/

import (
	"github.com/afoley587/52-weeks-of-projects/07-golang-gin-langchain/routes"
	"github.com/gin-gonic/gin"
)

/*
First, we instantiate a new application and we just
attach the router to it. All of the `/vacation/` endpoints
will be automatically added to our app.
*/

func main() {
	r := gin.Default()
	routes.GetVacationRouter(r)
	r.Run()
}
