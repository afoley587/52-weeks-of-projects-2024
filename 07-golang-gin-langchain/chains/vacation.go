package chains

import "github.com/google/uuid"

type Vacation struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}

var Vacations []Vacation

func GetVacationFromDb(id uuid.UUID) (Vacation, error) {
	return Vacation{}, nil
}

func GeneateVacationIdeaChange() {}
