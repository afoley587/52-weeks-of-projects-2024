from pydantic import BaseModel, Field
from typing import List


class InfluxWaveRecord(BaseModel):
    location: str = Field(description="Contents of the requested bucket")
    height: float = Field(description="Contents of the requested bucket")


class GetBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")


class ListBucketResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    records: List[InfluxWaveRecord] = Field(
        description="Contents of the requested bucket"
    )


class InsertWaveHeightRequest(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    location: str = Field(description="Contents of the requested bucket")
    height: float = Field(description="Contents of the requested bucket")


class InsertWaveHeightResponse(BaseModel):
    bucket: str = Field(description="Name of the requested bucket")
    location: str = Field(description="Contents of the requested bucket")
    height: float = Field(description="Contents of the requested bucket")