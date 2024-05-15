import os.path
from pathlib import Path
from pprint import pprint
from typing import Optional
from metadata.config import Config
from metadata.database.documents import get_by_name
from metadata.extractors.extractor import Extractor, FileMetadata
from typing_extensions import Self

from metadata.extractors.pdf_extractor import PdfExtractor


class FileExtractorBuilder:
    def __init__(self, builder: Optional[Extractor] = None):
        self._builder = builder
        self._config = Config()

    @staticmethod
    def get_file_extension(file_name: str) -> str:
        return Path(file_name).suffix

    @classmethod
    def prepare_file(cls, file_name: str) -> Self:
        extension = cls.get_file_extension(file_name)
        # TODO: Implement some kind of factory
        extractor: Optional[Extractor] = None

        document = get_by_name(file_name)
        if document is None:
            raise Exception("File not found")

        if extension.lower() == ".pdf":
            extractor = PdfExtractor(document)
        return FileExtractorBuilder(extractor)

    def process_file(self) -> FileMetadata:
        meta = self._builder.extract_meta_and_content()
        for m in meta:
            pprint(m)
