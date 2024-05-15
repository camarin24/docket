from sqlalchemy.orm import declarative_base, sessionmaker
from sqlalchemy import create_engine, Column, String

Base = declarative_base()

engine = create_engine('postgresql://postgres:postgres@localhost:5432/docket', echo=True)
Session = sessionmaker(engine)
