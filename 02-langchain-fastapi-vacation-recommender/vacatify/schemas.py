from pydantic import BaseModel, Field
from typing import List
import uuid


class GenerateVacationIdeaResponse(BaseModel):
    id: uuid.UUID = Field(description="ID Of the generated idea")
    completed: bool = Field(
        description="Flag indicating if the generation was completed"
    )


class GetVacationIdeaResponse(GenerateVacationIdeaResponse):
    idea: str = Field(description="The generated idea")


class GenerateVacationIdeaRequest(BaseModel):
    favorite_season: str = Field(description="Your favorite season")
    hobbies: List[str] = Field(description="The hobbies you enjoy")
    budget: int = Field(description="The budget for your vacation")
