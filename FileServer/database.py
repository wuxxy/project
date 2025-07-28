import datetime

from peewee import *
from dotenv import load_dotenv
load_dotenv()
import os


db = PostgresqlDatabase(os.getenv("DB_NAME"), user=os.getenv("DB_USER"), password=os.getenv("DB_PASS"), host=os.getenv("DB_HOST"), port=os.getenv("DB_PORT"))

class BaseModel(Model):
    """A base model that will use our Postgresql database"""
    class Meta:
        database = db

class Service(BaseModel):
    name = CharField(unique=True)
    secret = CharField(unique=True)
class File(BaseModel):
    name = CharField(unique=True)
    type = CharField(unique=False)
    location = CharField(unique=False)
    created = DateTimeField(default=datetime.datetime.now)
    modified = DateTimeField(default=datetime.datetime.now)
    parent_service = ForeignKeyField(Service)
def create_tables():
    with db:
        db.create_tables([Service, File])
