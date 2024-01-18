from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS
from pydantic import SecretStr
from typing import Union, List


class InfluxClient:
    def __init__(self, bucket: str, token: SecretStr, org: str, url: str) -> None:
        self.bucket = bucket
        self._client = InfluxDBClient(url=url, token=token.get_secret_value(), org=org)

    async def record_wave_height(self, location, height) -> None:
        p = Point("surf_heights").tag("location", location).field("height", height)
        await self._insert(p)
        return p

    async def read_wave_height(self, location: str = "", min_height: float = -1.0):
        query = f'from(bucket:"{self.bucket}")\
            |> range(start: -10m)\
            |> filter(fn:(r) => r._measurement == "surf_heights")'
        if location:
            query += f'|> filter(fn:(r) => r.location == "{location}")'
        if min_height > 0:
            query += f'|> filter(fn:(r) => r._field >= "{min_height}")'
        print(query)
        return await self._query(query)

    async def list_wave_heights(self):
        query = f'from(bucket:"{self.bucket}")\
            |> range(start: -10m)\
            |> filter(fn: (r) => r._measurement == "surf_heights")'
        return await self._query(query=query)

    async def _insert(self, p: Point) -> None:
        write_api = self._client.write_api(write_options=SYNCHRONOUS)
        write_api.write(bucket=self.bucket, record=p)

    async def _query(self, query: str = ""):
        query_api = self._client.query_api()
        result = query_api.query(query=query)
        results = []
        for table in result:
            for record in table.records:
                results.append((record.get_field(), record.get_value()))
        return results
