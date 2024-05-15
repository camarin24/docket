from pydantic_settings import BaseSettings


class Config(BaseSettings):
    docket_files_path: str


settings = Config()
