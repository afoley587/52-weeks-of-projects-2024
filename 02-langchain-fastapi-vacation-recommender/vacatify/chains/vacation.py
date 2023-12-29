import uuid
from typing import List

from langchain.chat_models import ChatOpenAI
from langchain.prompts import (
    ChatPromptTemplate,
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
)
from loguru import logger

from vacatify.schemas import Vacation

"""
`vacations` is our vacation "database". As previously noted, I say database in
quotes because this is just a dictionary that is shared across the system.
Ideally, this would be some more persistent/stable/scalable form of storage
but, for the purpose of this conversation, a dictionary is perfect.
"""
vacations = {}


"""
`generate_vacation_idea_chain` is where we finally start to invoke langchain.
It takes a few parameters:

1. The UUID which was passed in from our router. We will use this to save the
    results in our vacation database.
2. The users preferred season. We will use that as a parameter to the langchain chain.
3. The users favorite hobbies. We will use that as a parameter to the langchain chain.
4. The users financial budget. We will use that as a parameter to the langchain chain.

First, we create a system template and system message to pass to the LLM. A 
A system message is an instruction or information provided by the application or 
system to guide the conversation. The system message helps set the context and
instructions for the LLM and will guide how it responds to the human prompt.
A system template is just a templated form of the message.

A human message and template are the same idea. 

We can think of this like a chat application. The system prompt helps set 
up the chatbot. The human prompt is what the user would ask it.

Now that the templates are established, we can create a prompt from them
using the `from_template` methods. Next, we begin to intiialize the 
prompt template from the system message and the human message.

We can think of this as setting the scene for our chatbot: we gave them
a generic system message, we gave them a generic human message, and now
we can ask the LLM to respond to the prompts. 

To accomplish that, we have to use the `from_messages` method to begin
our chat conversation and then use `format_prompt` so that the prompt
gets formatted into text that the LLM will understand and that contains
all of the required context. 

Finally, we can chall our chain with `chat(request)` which submits the
formatted chat prompt to the LLM. When the LLM is done responding, we
can update our vacation database.

By this time, the user can then begin to query and read the response from the LLM
over the HTTP API.
"""


async def generate_vacation_idea_chain(
    id: uuid.UUID, season: str, hobbies: List[str], budget: int
):
    logger.info(f"idea generation starting for {id}")
    chat = ChatOpenAI()
    system_template = """
    You are an AI travel agent that will help me create a vacation idea.
    
    My favorite season is {season}.

    My hobbies include {hobbies}.

    My budget is {budget} dollars.
    """
    vacations[id] = Vacation(id=id, completed=False, idea="")

    system_message_prompt = SystemMessagePromptTemplate.from_template(system_template)
    human_template = "{travel_request}"
    human_message_prompt = HumanMessagePromptTemplate.from_template(human_template)
    chat_prompt = ChatPromptTemplate.from_messages(
        [system_message_prompt, human_message_prompt]
    )
    request = chat_prompt.format_prompt(
        season=season,
        budget=budget,
        hobbies=hobbies,
        travel_request="write a travel itinerary for me",
    ).to_messages()
    result = chat(request)
    vacations[id].idea = result.content
    logger.info(f"Completed idea generation for {id}")
