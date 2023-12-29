from pydantic import BaseModel, Field
from typing import List
import uuid


"""
The GenerateVacationIdeaResponse is what
we will return to a user that says the 
idea is currently being generated. Langchain might take
some time to generate the response, and we don't want users
to have to wait forever for their HTTP call to return. Because
of this, we will use background tasks (more on that later!)
and users can check in to see if their idea is finished after a 
few seconds.

The GenerateVacationIdeaResponse reflects this with two fields:

1. An ID field which will allow them to query our API for UUID of the
    project
2. A completed field which tells the user whether the idea generation
    is finished or not.
"""
class GenerateVacationIdeaResponse(BaseModel):
    id: uuid.UUID = Field(description="ID Of the generated idea")
    completed: bool = Field(
        description="Flag indicating if the generation was completed"
    )

"""
The GetVacationIdeaResponse is what we will return to a 
user when they query for the idea or its status. After
a few seconds, the user will say "Hm, is the idea done yet?"
and can query our API.
The GetVacationIdeaResponse has the same fields as GenerateVacationIdeaResponse,
but adds an idea field which is what the LLM will fill out when 
the generation is completed.
"""
class GetVacationIdeaResponse(GenerateVacationIdeaResponse):
    idea: str = Field(description="The generated idea")


"""
The GenerateVacationIdeaRequest is what a user will
provide to us so we can create a vacation idea for them.
We will expect them to tell us their favorite season,
any hobbies they may have, and what their vacation budget is.
We can feed these in to the LLM down the line.
"""
class GenerateVacationIdeaRequest(BaseModel):
    favorite_season: str = Field(description="Your favorite season")
    hobbies: List[str] = Field(description="The hobbies you enjoy")
    budget: int = Field(description="The budget for your vacation")


"""
The Vacation object will more or less be a data object for us.
It is identical to GetVacationIdeaResponse, but I sometimes
find it useful to have separate models for separate portions
of the stack so code is easier to maintain/modify later.
"""
class Vacation(GenerateVacationIdeaResponse):
    idea: str = Field(description="The generated idea")
