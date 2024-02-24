package routes

import "github.com/google/uuid"

/*
It's hard to talk about an API without understanding
what the request and response bodies are. Our API
will have one response body to support the POST
endpoint and then two response bodies, one to
support responses from POST endpoint and one
to support responses from the GET endpoint.

The `GenerateVacationIdeaRequest` is what a user will
provide to us so we can create a vacation idea for them.
We will expect them to tell us their favorite season,
any hobbies they may have, and what their vacation budget is.
We can feed these in to the LLM down the line.

The `GenerateVacationIdeaResponse` is what
we will return to a user that says the
idea is currently being generated. Langchain might take
some time to generate the response, and we don't want users
to have to wait forever for their HTTP call to return. Because
of this, we will use goroutines (more on that later!)
and users can check in to see if their idea is finished after a
few seconds.

The `GenerateVacationIdeaResponse` reflects this with two fields:

 1. An ID field which will allow them to query our API for UUID of the
    project
 2. A completed field which tells the user whether the idea generation
    is finished or not.

The `GetVacationIdeaResponse` is what we will return to a
user when they query for the idea or its status. After
a few seconds, the user will say "Hm, is the idea done yet?"
and can query our API.
The `GetVacationIdeaResponse` has the same fields as `GenerateVacationIdeaResponse`,
but adds an idea field which is what the LLM will fill out when
the generation is completed.
*/
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
