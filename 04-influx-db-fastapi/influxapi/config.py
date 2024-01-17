from pydantic_settings import BaseSettings
from pydantic import SecretStr
import os


class Settings(BaseSettings):
    influx_url: str = os.environ.get("INFLUX_URL")
    influx_token: SecretStr = os.environ.get("INFLUX_TOKEN")
    influx_org: str = os.environ.get("INFLUX_ORG")


settings = Settings()
