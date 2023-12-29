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

## Part I - Building the API
## Part II - Building the Chain
## Part III - Running And Testing