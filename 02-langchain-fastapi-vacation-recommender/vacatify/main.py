"""
Running

poetry run uvicorn vacatify.main:app --reload

Example Curls:

curl -X POST -H"Content-type: application/json" -d'{"favorite_season": "summer", "hobbies": ["this","that"], "budget":10}' http://localhost:8000/vacation/
curl -X GET -H"Content-type: application/json" http://localhost:8000/vacation/83a85995-fd55-4eeb-8526-c23c2468a641
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
