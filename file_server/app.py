import re
import magic
from flask import Flask, jsonify, request
from peewee import DatabaseError, IntegrityError
from manager import upload_file

TRAVERSAL_REGEX = re.compile(r'(^|[\\/])(\.\.?)([\\/]|$)')
ALLOWED_FILE_TYPES = ["png", "jpg", "jpeg", "gif", "webp", "mp4", "mp3", "wav", "avi", "mkv", "txt", "json", "xml"]
MIME_ALLOWLIST = {
    "image/png", "image/jpeg", "image/gif", "image/webp",
    "video/mp4", "video/x-matroska", "video/x-msvideo",
    "audio/mpeg", "audio/wav",
    "text/plain", "application/json", "application/xml", "text/xml"
}

mime_detector = magic.Magic(mime=True)

def detect_mime(file_stream) -> str:
    sample = file_stream.read(2048)
    file_stream.seek(0)
    return mime_detector.from_buffer(sample)
app = Flask(__name__)

@app.route("/upload", methods=["POST"])
def upload():
    file_name = request.form.get("name")
    file_service = request.form.get("ps_id")

    if 'file' not in request.files:
        return jsonify(error="No file provided"), 400
    if not file_name:
        return jsonify(error="No name provided"), 400
    if not file_service:
        return jsonify(error="No ps_id provided"), 400
    if TRAVERSAL_REGEX.search(file_name):
        return jsonify(error="Traversal is not allowed"), 400

    file = request.files['file']

    ext = file.filename.rsplit(".", 1)[-1].lower()
    if ext not in ALLOWED_FILE_TYPES:
        return jsonify(error="File extension not allowed"), 400

    mime = detect_mime(file.stream)
    if mime not in MIME_ALLOWLIST:
        return jsonify(error=f"Invalid MIME type: {mime}"), 400

    try:
        new_file = upload_file(file=file, file_name=file_name, file_type=ext, parent_service=file_service)
        return jsonify(success=True, file_id=new_file.id), 201
    except (IntegrityError, DatabaseError):
        return jsonify(error="Couldn't create file"), 500
    except Exception as e:
        print("Error uploading file:", e)
        return jsonify(error="An error occurred while uploading the file"), 500

if __name__ == "__main__":
    app.run(debug=True)
