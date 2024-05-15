import uuid

from metadata.database.base import Base, Session
from sqlalchemy import Column, Integer, String, UUID


class Document(Base):
    __tablename__ = "documents"
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    name = Column(String)
    storage_key = Column(String)
    original_path = Column(String)

    def __repr__(self):
        return f"<Document(id={self.id} ,name={self.name}, original_path={self.original_path})>"


def insert(doc: Document):
    with Session() as session:
        session.add(doc)
        session.commit()


def get_by_name(name: str) -> Document:
    with Session() as session:
        return session.query(Document).filter_by(name=name).first()
