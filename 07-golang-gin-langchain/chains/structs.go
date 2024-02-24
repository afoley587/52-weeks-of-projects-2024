package chains

import "github.com/google/uuid"

/*
Our vacation struct will just be data holder.
It holds the in-process and final vacation objects.
It has the same fields as the `GetVacationIdeaResponse`
we talked about earlier but I prefer to separate them
so it's easier to decouple these pieces of code.
*/
type Vacation struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}
