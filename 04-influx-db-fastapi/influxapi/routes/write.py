from fastapi import APIRouter, Request, BackgroundTasks, HTTPException

from influxapi.schemas import InsertBucketRequest, InsertBucketResponse
from influxapi.client.influx import InfluxClient
from influxapi.config import settings

write_router = APIRouter(prefix="/write")


@write_router.post(
    "/{bucket}/insert",
    summary="Insert data into a bucket.",
    responses={
        200: {"description": "Successfully Inserted Into Bucket."},
        400: {"description": "Bad data requested."},
        404: {"description": "Bucket not found."},
    },
)
async def insert_bucket(r: InsertBucketRequest, bucket: str) -> InsertBucketResponse:
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    ic.record_get_request()
    return InsertBucketResponse()
