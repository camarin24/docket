import os.path
from abc import ABC, abstractmethod
from typing import List, Optional

from pydantic import BaseModel
from metadata.config import settings


class FileMetadata(BaseModel):
    content: Optional[str]
    page_number: Optional[int]
    total_pages: Optional[int]
    title: Optional[str]
    author: Optional[str]
    subject: Optional[str]
    keywords: Optional[str]
    creator: Optional[str]
    creationDate: Optional[str]
    modDate: Optional[str]


class Extractor(ABC):
    def __init__(self, file_path: str):
        self.file_path = file_path
        # Thumbnail extractor
        # Raw contents extractor
        # Metadata extractor

    @abstractmethod
    def extract_thumbnail(self):
        # TODO: Validate if there is a metadata file already created
        pass

    @abstractmethod
    def extract_meta_and_content(self) -> List[FileMetadata]:
        # TODO: Validate if there is a metadata file already created
        pass

    @property
    def full_file_path(self) -> str:
        return os.path.join(settings.docket_files_path, self.file_path)
