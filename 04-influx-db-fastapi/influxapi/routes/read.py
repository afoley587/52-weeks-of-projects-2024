from fastapi import APIRouter, Request, HTTPException
from loguru import logger

from influxapi.schemas import GetBucketResponse, ListBucketResponse
from influxapi.client.influx import InfluxClient
from influxapi.config import settings

read_router = APIRouter(prefix="/read")


@read_router.get(
    "/{bucket}/list",
    summary="List's a bucket's contents.",
    responses={
        200: {"description": "Successfully Listed Bucket."},
        404: {"description": "Bucket not found."},
    },
)
async def list_bucket(r: Request, bucket: str) -> ListBucketResponse:
    logger.debug(f"Listing {bucket=}")
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    records = await ic.list_wave_heights()
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
async def query_bucket(
    r: Request, bucket: str, location: str = "", min_height: float = -1.0
) -> ListBucketResponse:
    logger.debug(f"Querying {bucket=} with {location=} and {min_height}")
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    records = await ic.read_wave_height(location=location, min_height=min_height)
    logger.debug(f"Records fetched {records=}")
    return ListBucketResponse(bucket=bucket, records=records)
