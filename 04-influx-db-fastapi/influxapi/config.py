from pydantic_settings import BaseSettings
from pydantic import SecretStr
import os

"""
All of the pieces of the puzzle are now built! We can then create a 
reusable set of settings that our clients can use. For this, we will
use pydantic's `BaseSettings` class. We will have three settings:

1. influx_url - The InfluxDB connection URL
2. influx_token - The InfluxDB authentication token
3. influx_org - The InfluxDB organization

These should look familiar from our routers!

"""


class Settings(BaseSettings):
    influx_url: str = os.environ.get("INFLUX_URL")
    influx_token: SecretStr = os.environ.get("INFLUX_TOKEN")
    influx_org: str = os.environ.get("INFLUX_ORG")


settings = Settings()
