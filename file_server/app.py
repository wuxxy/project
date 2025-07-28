import re

from flask import Flask, jsonify, request
from peewee import DatabaseError, IntegrityError

from manager import upload_file
TRAVERSAL_REGEX = re.compile(r'(^|[\\/])(\.\.?)([\\/]|$)')

ALLOWED_FILE_TYPES = ["png", "jpg", "jpeg", "gif", "pdf", "txt", "docx",
                      "xlsx", "csv", "json", "xml", "webp", "mp4", "mp3",
                      "wav", "avi", "mkv", "zip", "rar"]
app = Flask(__name__)
@app.route("/upload", methods=["POST"])
def upload():
    file_name = request.form.get("name")
    file_location = request.form.get("location")
    file_type = request.form.get("type")
    file_service = request.form.get("type")
    if 'file' not in request.files:
        return jsonify(error="No file provided"), 400
    if file_name is None or file_name == "":
        return jsonify(error="No name provided"), 400
    if file_type is None or file_type == "":
        return jsonify(error="No type provided"), 400
    if TRAVERSAL_REGEX.search(file_name):
        return jsonify(error="Traversal is not allowed"), 400
    file = request.files['file']
    if file_type not in ALLOWED_FILE_TYPES:
        return jsonify(error="File type not allowed"), 400
    try:
        # Step 1: Create and save the DB record
        new_file = upload_file(file=file, file_name=file_name, file_type=file_type)
        return jsonify(success=True, file_id=new_file.id), 201
    except IntegrityError(DatabaseError):
        return jsonify(error="Couldn't create file")
    except Exception as e:
        print("Error uploading file:", e)
        return jsonify(error="An error occurred while uploading the file"), 500
