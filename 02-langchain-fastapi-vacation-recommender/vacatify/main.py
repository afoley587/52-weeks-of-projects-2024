from fastapi import FastAPI

from vacatify.routes.vacation import vacation_router

app = FastAPI()
app.include_router(vacation_router)
