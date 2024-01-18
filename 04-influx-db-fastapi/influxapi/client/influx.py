from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS
from influxdb_client.rest import ApiException
from pydantic import SecretStr
from typing import List, Any
from loguru import logger
from urllib3.exceptions import NewConnectionError
from influxapi.schemas import InfluxWaveRecord

"""
Now, we need to implement a way to interact with our Influx database.
We also want to handle some exceptions so that our API doesn't return
error 500 codes on unhandled exceptions. We have three expections that we
expect to handle. The first being the `InfluxNotAvailableException` which
is what will be raised when InfluxDB can't be reached. Next, we have the
`BucketNotFoundException` which is what will be raised if a user requests
a bucket doesn't exist. The last being the `BadQueryException` which
will get raised if there's an error in our queries.

With our exceptions out of the way, we can build our InfluxDB interface.
The `InfluxWaveClient` will be initialized with a bucket, a token, an
organization, and a url. The bucket will be used when reading/inserting
data into InfluxDB. The URL, token, and organization will be used to connect
to the right InfluxDB instance. The client provides a few "public" methods
to users: `record_wave_height`, `read_wave_height`, and `list_wave_heights`.
It also has two "private" methods: `_insert` and `_query`.

First, we will discuss our "public" methods. `record_wave_height` takes a
few parameters from the caller: a location to record and the wave's height
to record. It create a Point object and then calls the "private" `_insert`
method with that point. Next, we have the `read_wave_height` method. This
method also takes two parameters: a location to filter for and a minimmum 
height to filter on. For example, if we pass "hawaii" and "1.25", we would
be looking for waves in Hawaii that are at least 1.25 (unit doesn't matter).
This would call the "private" `_query` method with the relevant filters
and return the matching data points to the caller. The `list_wave_heights`
method does almost the same thing. It just calls the `read_wave_height` 
method with the default/empty parameters which would match all data points
in the database.

The "private" methods are `_insert` and `_query`. `_insert` will take a 
data point from the caller. It will use InfluxDB's `write_api` to store
the data point in the database. The `_query` method uses InfluxDB's
`query_api` to send the given query to the database. It then puts all
of the records returned from the `query_api` into the pydantic model 
we discussed above in step 1.
"""


class InfluxNotAvailableException(Exception):
    STATUS_CODE = 503
    DESCRIPTION = "Unable to connect to influx."


class BucketNotFoundException(Exception):
    STATUS_CODE = 404
    DESCRIPTION = "Bucket Not Found."


class BadQueryException(Exception):
    STATUS_CODE = 400
    DESCRIPTION = "Bad Query."


class InfluxWaveClient:
    """A restricted client which implements an interface
    to query the wave-related data from the Influx database
    """

    MEASUREMENT_NAME: str = "surf_heights"

    def __init__(self, bucket: str, token: SecretStr, org: str, url: str) -> None:
        self.bucket = bucket
        self._client = InfluxDBClient(url=url, token=token.get_secret_value(), org=org)

    async def record_wave_height(self, location: str, height: float) -> None:
        """Records a new wave height for a given location

        Arguments:
            location (str): The location to tag the data point as
            height (float): The height of the measured wave

        Returns:
            None
        """
        location = location.lower()
        p = (
            Point(InfluxWaveClient.MEASUREMENT_NAME)
            .tag("location", location)
            .field("height", height)
        )
        await self._insert(p)

    async def read_wave_height(
        self, location: str = "", min_height: float = -1.0
    ) -> List[InfluxWaveRecord]:
        """Reads a wave height given a specific

        Arguments:
            location (str): The location to filter on
            min_height (float): The minimum wave height to filter on

        Returns:
            res (List[InfluxWaveRecord]): The datapoints that match this filter
        """
        query = f'from(bucket:"{self.bucket}")\
            |> range(start: -10m) \
            |> filter(fn:(r) => r._measurement == "{InfluxWaveClient.MEASUREMENT_NAME}")'
        if location:
            location = location.lower()
            query += f'|> filter(fn:(r) => r.location == "{location}")'
        if min_height > 0:
            query += f'|> filter(fn:(r) => r._field >= "{min_height}")'
        return await self._query(query)

    async def list_wave_heights(self) -> List[InfluxWaveRecord]:
        """Lists the bucket in question

        Arguments:
            None

        Returns:
            res (List[InfluxWaveRecord]): All waves in the buckets
        """
        return await self.read_wave_height(location="", min_height=-1.0)

    async def _insert(self, p: Point) -> Any:
        """Inserts a point into the database via InfluxDB write_api

        Arguments:
            p (Point): The data point to insert into the database

        Returns:
            res (Any): Results from the write_api
        """
        write_api = self._client.write_api(write_options=SYNCHRONOUS)
        try:
            res = write_api.write(bucket=self.bucket, record=p)
        except NewConnectionError:
            raise InfluxNotAvailableException()
        except ApiException as e:
            if e.status and e.status == 404:
                raise BucketNotFoundException()
            raise InfluxNotAvailableException()
        logger.info(f"{res=}")
        return res

    async def _query(self, query: str = "") -> List[InfluxWaveRecord]:
        """Queries the InfluxDB with the provided query string

        Arguments:
            query (str): The raw query string to pass to InfluxSB

        Returns:
            res (List[InfluxWaveRecord]): A list of waves that match the query
        """
        logger.debug(f"Running {query=}")
        query_api = self._client.query_api()
        try:
            result = query_api.query(query=query)
        except NewConnectionError:
            raise InfluxNotAvailableException()
        except ApiException as e:
            if e.status and e.status == 404:
                raise BadQueryException()
            if e.status and e.status == 404:
                raise BucketNotFoundException()
            raise InfluxNotAvailableException()
        res = []
        for table in result:
            for record in table.records:
                r = InfluxWaveRecord(
                    location=record.values.get("location"), height=record.get_value()
                )
                res.append(r)
        logger.debug(f"Query returned {len(res)} records")
        return res
