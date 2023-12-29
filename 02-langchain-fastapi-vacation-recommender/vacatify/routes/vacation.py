import uuid

from fastapi import APIRouter, Request, BackgroundTasks, HTTPException

from vacatify.schemas import (
    GenerateVacationIdeaResponse,
    GetVacationIdeaResponse,
    GenerateVacationIdeaRequest,
)

from vacatify.chains.vacation import generate_vacation_idea_chain, vacations

"""
The `vacation_router` is the router that will be used
by the main application. All of the endpoints on the router
will therefore be added to the application so we can 
begin to leverage them!
"""
vacation_router = APIRouter(prefix="/vacation")


"""
The first endpoint we will add is a POST endpoint at
`http://the-api-ip:the-api-port/vacation/`. The main purposes
of this endpoint are to:

1. Start a new background task to create a vacation idea for you
2. Report the created ID and progress back to the user

We can also see that this endpoint takes, as input, a request
of type `GenerateVacationIdeaRequest` and returns a response
of type `GenerateVacationIdeaResponse`. These data types were
defined above in our schemas section. Thankfully, fastAPI does
all of the required serialization both to and from the API!

This method also leverages a FastAPI notion called BackgroundTasks.
We can define background tasks to be run after returning a response.
This is useful because, at time, Langchain might take a while to run.
You don't want your clients to have to handle huge timeouts beacuse that
might indicate other issues. By scheduling these as background tasks, we can
say "hey client, the work has started. Check back later :)"

So, let's look more in depth at `generate_vacation()`. First, 
we generate a UUID for the idea. Then, we submit the chain
with all of the parameters from our user to our background task.
The observant reader will see the function `generate_vacation_idea_chain`.
This is the function that will run langchain and update our vacation
database. We will talk in more depth about it in a following section.
Finally, we return the response with the ID and the completed 
flag set to false.
"""


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


"""
Our second endpoint is a GET endpoint at
`http://the-api-ip:the-api-port/vacation/<id>`. The main
purpose of this endpoint is going to be use to either 
poll/query/read the idea that was created from the generation above.

We can see that is accepts a UUID id as a parameter and returns a 
`GetVacationIdeaResponse` which was again defined in our schema section
above.

We can see that this endpoint just looks to see if a vacation ID matching
the requested ID is present in our vacation "database". I say database in
quotes because this is just a dictionary that is shared across the system.
Ideally, this would be some more persistent/stable/scalable form of storage
but, for the purpose of this conversation, a dictionary is perfect.

If we know the ID that the user is requesting, we can put all the relevant
data into a `GetVacationIdeaResponse` and return it to the user. Otherwise,
we throw a 404.
"""


@vacation_router.get(
    "/{id}",
    summary="Get the generated a vacation idea.",
    responses={
        200: {"description": "Successfully fetched vacation."},
        404: {"description": "Vacation not found."},
    },
)
async def get_vacation(r: Request, id: uuid.UUID) -> GetVacationIdeaResponse:
    """Returns the vacation generation for you."""
    if id in vacations:
        vacay = vacations[id]
        return GetVacationIdeaResponse(
            id=vacay.id, completed=vacay.completed, idea=vacay.idea
        )
    raise HTTPException(status_code=404, detail="ID not found")
