"""
Running

poetry run uvicorn vacatify.main:app --reload

Example Curls:

curl -X POST -H"Content-type: application/json" -d'{"favorite_season": "summer", "hobbies": ["surfing","running"], "budget":1000}' http://localhost:8000/vacation/
curl -X GET -H"Content-type: application/json" http://localhost:8000/vacation/cfc8c891-6826-4320-a652-bd6febd9fd7b
"""
from fastapi import FastAPI

from vacatify.routes.vacation import vacation_router

"""
First, we instantiate a new application and we just
attach the router to it. All of the `/vacation/` endpoints
will be automatically added to our app.
"""
app = FastAPI()
app.include_router(vacation_router)
