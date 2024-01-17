from fastapi import APIRouter, Request, HTTPException
from influxapi.schemas import GetBucketResponse, ListBucketResponse
from influxapi.client.influx import InfluxClient
from influxapi.config import settings

read_router = APIRouter(prefix="/read")


@read_router.get(
    "/{bucket}",
    summary="Gets a bucket's metadata.",
    responses={
        200: {"description": "Successfully Found Bucket."},
        404: {"description": "Bucket not found."},
    },
)
async def get_bucket(r: Request, bucket: str) -> GetBucketResponse:
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    data = ic.get_bucket()
    return GetBucketResponse(bucket=bucket)


@read_router.get(
    "/{bucket}/list",
    summary="List's a bucket's contents.",
    responses={
        200: {"description": "Successfully Listed Bucket."},
        404: {"description": "Bucket not found."},
    },
)
async def list_bucket(r: Request, bucket: str) -> ListBucketResponse:
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    records = ic.list_bucket()
    return ListBucketResponse(bucket=bucket, records=records)


@read_router.get(
    "/{bucket}/query",
    summary="Queries a bucket's contents.",
    responses={
        200: {"description": "Successfully Queried Bucket."},
        400: {"description": "Bad Filter Requested."},
        404: {"description": "Bucket not found."},
    },
)
async def query_bucket(r: Request, bucket: str) -> ListBucketResponse:
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    records = ic.query_bucket()
    return ListBucketResponse(bucket=bucket, records=records)
