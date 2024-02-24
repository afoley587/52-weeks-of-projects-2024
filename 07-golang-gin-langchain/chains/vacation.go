package chains

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"golang.org/x/exp/slices"
)

/*
The scene has been set with the API. Now, we need a way
to converse with our LLM (or at least as it questions).

Let's define a "database" to store all of the generated ideas.
`Vacations` is our vacation "database". I say database in
quotes because this is just a map that is shared across this package.
Ideally, this would be some more persistent/stable/scalable form of storage
but, for the purpose of this conversation, a map is perfect.

We need to provide two methods to whoever wants to use this package:

1. A way for the caller to retrieve a vacation from our "dataabase"
2. A way for the caller to request a new vacation idea to be generated

To tackle number 1, we will write the `GetVacationFromDb(id uuid.UUID)`
function. This function will take the ID of the vacation. It then tries
to find the vacation in the map and, if it exists, it returns the vacation
object. Otherwise, it returns an error if the ID does not exist in the
database.

Next, we need something to actually create ideas and store them into our
database.

`GeneateVacationIdeaChange` is where we finally start to invoke langchain.
It takes a few parameters:

1. The UUID which was passed in from our router. We will use this to save the
    results in our vacation database.
2. The users preferred season. We will use that as a parameter to the langchain chain.
3. The users favorite hobbies. We will use that as a parameter to the langchain chain.
4. The users financial budget. We will use that as a parameter to the langchain chain.

First, we need to instantiate our LLM model (openai in this case). Then
we need to create a few prompts. We create a system prompt to pass to the LLM. A
A system prompt is an instruction or information provided by the application or
system to guide the conversation. The system prompt helps set the context and
instructions for the LLM and will guide how it responds to the human prompt.

A human message and template follows the same idea. We can think of this like a
chat application. The system prompt helps set up the chatbot.
The human prompt is what the user would ask it.

Now that the templates are established, we can create a chat prompt from them
by first creating a chat prompt template. To do this, we use the `FormatMessages` method
to then insert our user-provided values into our templates. Now everything is
templated in string format. We will create LLM message contents which is what our
llm will expect as input. Finally, we can invoke our LLM with `GenerateContent`.
The output of `GenerateContent` will the result from the OpenAI API but, we only
care about the content that the LLM generated. The content is the string response
that would be returned to you in something like the ChatGPT window.
*/

var Vacations []*Vacation

func GetVacationFromDb(id uuid.UUID) (Vacation, error) {
	// Use the slices package to find the index of the object with
	// matching ID in the database. If it does not exist, this will return
	// -1
	idx := slices.IndexFunc(Vacations, func(v *Vacation) bool { return v.Id == id })

	// If the ID didn't exist, return an error and let the caller
	// handle it
	if idx < 0 {
		return Vacation{}, errors.New("ID Not Found")
	}

	// Otherwise, return the Vacation object
	return *Vacations[idx], nil
}

func GeneateVacationIdeaChange(id uuid.UUID, budget int, season string, hobbies []string) {
	log.Printf("Generating new vacation with ID: %s", id)

	// Create a new vacation object and add it to our database. Initially,
	// the idea field will be empty and the completed flag will be false
	v := &Vacation{Id: id, Completed: false, Idea: ""}
	Vacations = append(Vacations, v)

	// Create a new OpenAI LLM Object
	ctx := context.Background()
	llm, err := openai.New()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Create a system prompt with the season, hobbies, and budget parameters
	// Helps tell the LLM how to act / respond to queries
	system_message_prompt_string := "You are an AI travel agent that will help me create a vacation idea.\n" +
		"My favorite season is {{.season}}.\n" +
		"My hobbies include {{.hobbies}}.\n" +
		"My budget is {{.budget}} dollars.\n"
	system_message_prompt := prompts.NewSystemMessagePromptTemplate(system_message_prompt_string, []string{"season", "hobbies", "dollars"})

	// Create a human prompt with the request that a human would have
	human_message_prompt_string := "write a travel itinerary for me"
	human_message_prompt := prompts.NewHumanMessagePromptTemplate(human_message_prompt_string, []string{})

	// Create a chat prompt consisting of the system messages and human messages
	// At this point, we will also inject the values into the prompts
	// and turn them into message content objects which we can feed through
	// to our LLM
	chat_prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{system_message_prompt, human_message_prompt})

	vals := map[string]any{
		"season":  season,
		"budget":  budget,
		"hobbies": strings.Join(hobbies, ","),
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

	// Invoke the LLM with the messages which
	completion, err := llm.GenerateContent(ctx, content)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	v.Idea = completion.Choices[0].Content
	v.Completed = true

	log.Printf("Generation for %s is done!", v.Id)
}
