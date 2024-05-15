from typing import List

from metadata.database.documents import Document
from metadata.extractors.extractor import Extractor, FileMetadata
from unstructured.partition.pdf import partition_pdf


# TODO: Add unstructured as docker container

class PdfExtractor(Extractor):

    def __init__(self, document: Document):
        self.document = document
        super().__init__(document.name)

    def extract_thumbnail(self):
        pass

    def extract_meta_and_content(self) -> List[FileMetadata]:
        print(self.full_file_path)
        elements = partition_pdf(self.full_file_path, strategy="hi_res", infer_table_structure=True)
        return [dict(
            content=e.text,
            page_number=e.metadata.page_number,
            title=e.metadata.subject
        ) for e in elements]
