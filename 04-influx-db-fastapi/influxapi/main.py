"""
Running

poetry run uvicorn influxapi.main:app --reload
"""
from fastapi import FastAPI

from influxapi.routes.read import read_router
from influxapi.routes.write import write_router

"""
Finally, we can attach our routers to our FastAPI App and use uvicorn to kick 
it off.
"""
app = FastAPI()
app.include_router(read_router)
app.include_router(write_router)
