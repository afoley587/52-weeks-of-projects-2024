from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS
from pydantic import SecretStr
from typing import Union, List


class InfluxClient:
    def __init__(self, bucket: str, token: SecretStr, org: str, url: str) -> None:
        self.bucket = bucket
        self._client = InfluxDBClient(url=url, token=token.get_secret_value(), org=org)

    def record_get_request(self) -> None:
        self._insert()

    def read_get_request(self) -> List[Point]:
        self._query()

    def list_bucket(self):
        pass

    def query_bucket(self):
        pass

    def _insert(self, d: Point) -> None:
        write_api = self._client.write_api(write_options=SYNCHRONOUS)
        write_api.write(bucket=self.bucket, record=d)

    def _query(self) -> List[Point]:
        query_api = self._client.query_api()
        query = 'from(bucket:"my-bucket")\
        |> range(start: -10m)\
        |> filter(fn:(r) => r._measurement == "my_measurement")\
        |> filter(fn:(r) => r.location == "Prague")\
        |> filter(fn:(r) => r._field == "temperature")'
        result = query_api.query(query=query)
        results = []
        for table in result:
            for record in table.records:
                results.append((record.get_field(), record.get_value()))
