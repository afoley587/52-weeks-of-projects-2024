import uuid

from fastapi import APIRouter, Request, BackgroundTasks

from vacatify.schemas import (
    GenerateVacationIdeaResponse,
    GetVacationIdeaResponse,
    GenerateVacationIdeaRequest,
)

from vacatify.chains.vacation import generate_vacation_idea_chain, vacations

vacation_router = APIRouter(prefix="/vacation")

# curl -X POST -H"Content-type: application/json" -d'{"favorite_season": "summer", "hobbies": ["this","that"], "budget":10}' http://localhost:8000/vacation/
# curl -X GET -H"Content-type: application/json" http://localhost:8000/vacation/83a85995-fd55-4eeb-8526-c23c2468a641


@vacation_router.post(
    "/",
    summary="Generate a vacation idea.",
    responses={
        201: {"description": "Successfully initiated task."},
    },
)
async def generate_vacation(
    r: GenerateVacationIdeaRequest, background_tasks: BackgroundTasks
) -> GenerateVacationIdeaResponse:
    """Initiates a vacation generation for you."""

    idea_id = uuid.uuid4()
    background_tasks.add_task(
        generate_vacation_idea_chain,
        idea_id,
        r.favorite_season,
        r.hobbies,
        r.budget,
    )
    return GenerateVacationIdeaResponse(id=idea_id, completed=False)


@vacation_router.get(
    "/{id}",
    summary="Get the generated a vacation idea.",
    responses={
        201: {"description": "Successfully initiated task."},
    },
)
async def get_vacation(r: Request, id: uuid.UUID) -> GetVacationIdeaResponse:
    """Returns the vacation generation for you."""
    if id in vacations:
        return GetVacationIdeaResponse(id=id, completed=True, idea=vacations[id])
    return GetVacationIdeaResponse(id=id, completed=False, idea="str")
