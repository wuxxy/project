import datetime
from database import *

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

    # Step 2: Save the file to disk
    file_path = f"files/{file_name}_{new_file.id}"
    file.save(file_path)

    # Step 3: Update the file entry with actual location and size
    new_file.location = file_path
    new_file.file_size = file.content_length  # ← get real size
    new_file.save()
    return new_file