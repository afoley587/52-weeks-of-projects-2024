from pydantic import BaseModel, Field
from typing import List

"""
Schemas are what our APIs and Client will use to fetch data
and put them into some normalized, jsonifiable format. Typically, 
we would have two different sets of files: `schemas.py` and `models.py`.
Schemas would typically be used for REST API responses/requests while
models might be used closer to the data layer and deal with the database
result sets. However, in our case, we will just combine the two for 
brevity and because the size of our API is so small.

First, we have our `InfluxWaveRecord` model. This model is our only real
database model. When the client fetches data from InfluxDB, it will put
each record into a `InfluxWaveRecord` type and then return a list of those
to the caller. More on that later! The `InfluxWaveRecord` model has two
attributes: a location of the wave and the height of the wave.

Next, we have a request/response pair for when we want to insert data
into InfluxDB. These seem like duplicates, but I like to keep request
and response models separate so that updating one (i.e. providing
more information in a response) is easy and makes our code more flexible.

The `InsertWaveHeightRequest` will be used when a user sends data to
our API. It has two attributes: a location of the wave and the height 
of the wave. The `InsertWaveHeightResponse` is what will be returned
to the user after they insert data into InfluxDB. It has the same
two attributes as `InsertWaveHeightRequest`.

Finally, we have a response for when a user tries to read/list/query
InfluxDB. The `ListBucketResponse` has two attributes: the bucket
that was queried and a list of all of the `InfluxWaveRecord`
that were returned (either a listing of the entire bucket
or a filtered/queried subset).
"""


class InfluxWaveRecord(BaseModel):
    location: str = Field(description="Location of the recorded wave")
    height: float = Field(description="Height of the recorded wave")


class InsertWaveHeightRequest(BaseModel):
    location: str = Field(description="Location of the recorded wave")
    height: float = Field(description="Height of the recorded wave")


class InsertWaveHeightResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    location: str = Field(description="Location of the recorded wave")
    height: float = Field(description="Height of the recorded wave")


class ListBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    records: List[InfluxWaveRecord] = Field(
        description="Contents of the requested bucket"
    )
