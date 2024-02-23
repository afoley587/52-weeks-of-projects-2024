package chains

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"golang.org/x/exp/slices"
)

type Vacation struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}

var Vacations []*Vacation

func GetVacationFromDb(id uuid.UUID) (*Vacation, error) {
	idx := slices.IndexFunc(Vacations, func(v *Vacation) bool { return v.Id == id })

	if idx < 0 {
		return nil, errors.New("ID Not Found")
	}

	return Vacations[idx], nil
}

func GeneateVacationIdeaChange(id uuid.UUID) {
	v := &Vacation{Id: id, Completed: false, Idea: ""}
	Vacations = append(Vacations, v)

	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	stringPrompt := "You are an AI travel agent that will help me create a vacation idea.\n" +
		"My favorite season is {{.season}}.\n" +
		"My hobbies include {{.hobbies}}.\n" +
		"My budget is {{.budget}} dollars.\n"

	system_message_prompt := prompts.NewSystemMessagePromptTemplate(stringPrompt, []string{"season", "hobbies", "dollars"})

	human_message_prompt := prompts.NewHumanMessagePromptTemplate("write a travel itinerary for me", []string{})

	chat_prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{system_message_prompt, human_message_prompt})

	vals := map[string]any{
		"season":  "summer",
		"budget":  "10",
		"hobbies": "surfing",
	}
	msgs, err := chat_prompt.FormatMessages(vals)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	content := []llms.MessageContent{
		llms.TextParts(msgs[0].GetType(), msgs[0].GetContent()),
		llms.TextParts(msgs[1].GetType(), msgs[1].GetContent()),
	}

	/*
		completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
	*/
	completion, err := llm.GenerateContent(ctx, content)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	v.Idea = completion.Choices[0].Content
	v.Completed = true

	log.Printf("Generation for %s is done!", v.Id)
}
