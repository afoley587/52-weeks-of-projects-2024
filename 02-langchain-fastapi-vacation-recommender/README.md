# Hosting Langchain With FastAPI - Vacation Planner Project

Large language models, like ChatGPT, Llama, etc., have taken
the world by storm over the past 12 months. If you don't believe me,
let's take a look at how long it took ChatGPT to get to 1 million 
users as compared to some other behemoths. Graph is courtesy
of [exploding topics](https://explodingtopics.com/blog/chatgpt-users).

![chatgpt-1m-users](./images/chatgpt-1m-users.png)

Looking at this, we can see it took chatGPT 5 days to get to 1 million users.
Netflix took almost 4 years to get there. Now, I know that times have changed
and the internet era has boomed which helps these companies grow. But still,
5 days, REALLY?!

LLMs in general have become applicable to almost every single industry from 
health care to finance. If you're in any industry that has been impacted by
LLMs, it's probably a good idea to learn one or two things about them. If
you're like me, who writes software for a living, it's a good idea to know
how to use it and what some of the stuff means.

As always, breaking in to a new field is so difficult because you have to cut through
the technical jargon, find the right place to start, and find a way
to learn the new tech so that you don't feel obsolete. I don't know how you
learn best, but I am a learn-by-doer. I like to combine something I know
really well with something I don't. By doing that, I see where the new
puzzle piece fits into my overall puzzle (my brain).

That's exactly what I'm doing today. I've written hundreds of APIs in
python using almost every framework under the sun (Django, flask, tornado, you name it).
So, I'm going to combine that with something that's a bit newer to me - Langchain.
I've done a few courses in Langchain and I've watched some videos on it but, I'm still
learning, and I want to become production ready!

So, let's get started by introducing our stars of the show: FastAPI and Langchain.

## What Is FastAPI
FastAPI is a modern, fast (as the name implies), web framework for building APIs with 
Python. It is designed to be easy to use, efficient, and highly performant, 
leveraging the power of Python type hints for automatic data validation and 
documentation generation. Asynchronous programming is a core feature, allowing for 
high concurrency and scalability.

## What Is LangChain
Langchain, on the other hand, is a framework for developing applications powered by 
language models. It is designed to let you effortlessly plug in to an LLM and
enables users to provide extra context to the LLM. Simply put, LangChain enables LLM 
models to generate responses based on the most up-to-date information available online, 
in documents, or from other data sources.

## Installing
Both FastAPI and Langchain are python packages, so they can be pip installed:

```shell
prompt> pip install fastapi langchain
```

But there are a lot of supporting libraries that are useful in this project, so
I have included a [requirements.txt]() that you can get off and running with.

## Part 0 - File Structure

First, let's look at our file structure:
```shell
.
├── requirements.txt
└── vacatify
    ├── __init__.py
    ├── chains
    │   ├── __init__.py
    │   └── vacation.py
    ├── main.py
    ├── routes
    │   ├── __init__.py
    │   └── vacation.py
    └── schemas.py
```

Let's also give a high level overview of the files/directories in this tree. 
I will omit the `__init__.py` files from the description.

1. `vacatify`: The root level of our python package
2. `chains`: We will add our langchain logic in here. For example, we will generate
    and format prompts around vacations in the `vacation.py` file. We will then submit 
    these prompts to langchain, let it query the LLM, and save it's response.
3. `main.py`: This will be our application's entrypoint. It will start the server, 
    attach the API routes, etc.
4. `routes`: These will be the routes within our application. So, `vacation.py` will
    house the HTTP endpoints which pertain to adding a new vacation idea.
5. `schemas.py`: This file will house the request/response/pydantic schemas we want to 
    be able to use within our API.

## Part I - Building the API

### schemas.py
Let's start with the schemas so we know what our API is passing around.

The `GenerateVacationIdeaResponse` is what
we will return to a user that says the 
idea is currently being generated. Langchain might take
some time to generate the response, and we don't want users
to have to wait forever for their HTTP call to return. Because
of this, we will use background tasks (more on that later!)
and users can check in to see if their idea is finished after a 
few seconds.

The `GenerateVacationIdeaResponse` reflects this with two fields:

1. An ID field which will allow them to query our API for UUID of the
    project
2. A completed field which tells the user whether the idea generation
    is finished or not.

The `GetVacationIdeaResponse` is what we will return to a 
user when they query for the idea or its status. After
a few seconds, the user will say "Hm, is the idea done yet?"
and can query our API.
The `GetVacationIdeaResponse` has the same fields as `GenerateVacationIdeaResponse`,
but adds an idea field which is what the LLM will fill out when 
the generation is completed.

The `GenerateVacationIdeaRequest` is what a user will
provide to us so we can create a vacation idea for them.
We will expect them to tell us their favorite season,
any hobbies they may have, and what their vacation budget is.
We can feed these in to the LLM down the line.


The `Vacation` object will more or less be a data object for us.
It is identical to `GetVacationIdeaResponse`, but I sometimes
find it useful to have separate models for separate portions
of the stack so code is easier to maintain/modify later.

```python
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



class Vacation(GenerateVacationIdeaResponse):
    idea: str = Field(description="The generated idea")
```
## Part II - Building the Chain
## Part III - Running And Testing