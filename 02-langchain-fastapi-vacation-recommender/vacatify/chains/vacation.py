import uuid
from typing import List

from langchain.chat_models import ChatOpenAI
from langchain.prompts import (
    ChatPromptTemplate,
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
)
from langchain.schema import AIMessage, HumanMessage, SystemMessage
from loguru import logger

from vacatify.schemas import Vacation

vacations = {}


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
