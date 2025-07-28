import datetime
import uuid

from peewee import *
from dotenv import load_dotenv
load_dotenv()
import os


db = PostgresqlDatabase(os.getenv("DB_NAME"), user=os.getenv("DB_USER"), password=os.getenv("DB_PASS"), host=os.getenv("DB_HOST"), port=os.getenv("DB_PORT"))

class BaseModel(Model):
    created = DateTimeField(default=datetime.datetime.utcnow)
    modified = DateTimeField(default=datetime.datetime.utcnow)

    def save(self, *args, **kwargs):
        if not self.created:  # New record
            self.created = datetime.datetime.utcnow()
        self.modified = datetime.datetime.utcnow()
        return super().save(*args, **kwargs)
    class Meta:
        database = db

class Service(BaseModel):
    id = UUIDField(primary_key=True, default=uuid.uuid4)
    name = CharField(unique=True)
    secret = CharField(unique=True)
class File(BaseModel):
    id = UUIDField(primary_key=True, default=uuid.uuid4)
    name = CharField(unique=False)
    type = CharField(unique=False)
    file_size = IntegerField(default=0)
    parent_service = ForeignKeyField(Service)


def create_tables():
    with db:
        db.create_tables([Service, File])
        print("Tables created")
