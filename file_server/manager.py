import datetime
from database import *
import gzip
import io
def gzip_bytes(data: bytes) -> bytes:
    out = io.BytesIO()
    with gzip.GzipFile(fileobj=out, mode='wb') as f:
        f.write(data)
    return out.getvalue()

def upload_file(file, file_name, file_type, parent_service=None):
    new_file = File.create(
        name=file_name,
        location="",  # temp placeholder
        type=file_type,
        created=datetime.datetime.now(datetime.UTC),
        modified=datetime.datetime.now(datetime.UTC),
        parent_service=parent_service,
        file_size=0
    )
    file_bytes = file.read()
    file.seek(0)
    # Step 2: Save the file to disk
    file_path = f"files/{file_name}_{new_file.id}.gz"
    gzipped_data = gzip_bytes(file_bytes)
    with open(file_path, "wb") as f:
        f.write(gzipped_data)
    # Step 3: Update the file entry with actual location and size
    new_file.location = file_path
    new_file.file_size = file.content_length  # ← get real size
    new_file.save()
    return new_file