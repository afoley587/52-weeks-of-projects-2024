from pydantic import BaseModel, Field
from typing import List
import uuid


class GetBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")


class ListBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    records: List[str] = Field(description="Contents of the requested bucket")


class InsertBucketRequest(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    records: List[str] = Field(description="Contents of the requested bucket")


class InsertBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    records: List[str] = Field(description="Contents of the requested bucket")
