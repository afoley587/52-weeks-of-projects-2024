from fastapi import APIRouter, Request, BackgroundTasks, HTTPException

from influxapi.schemas import InsertWaveHeightRequest, InsertWaveHeightResponse
from influxapi.client.influx import InfluxClient
from influxapi.config import settings

write_router = APIRouter(prefix="/write")


@write_router.post(
    "/{bucket}/insert",
    summary="Insert data into a bucket.",
    responses={
        201: {"description": "Successfully Inserted Into Bucket."},
        400: {"description": "Bad data requested."},
        404: {"description": "Bucket not found."},
    },
)
async def insert_bucket(
    r: InsertWaveHeightRequest, bucket: str
) -> InsertWaveHeightResponse:
    ic = InfluxClient(
        bucket, settings.influx_token, settings.influx_org, settings.influx_url
    )
    await ic.record_wave_height(r.location, r.height)
    return InsertWaveHeightResponse(bucket=bucket, location=r.location, height=r.height)
