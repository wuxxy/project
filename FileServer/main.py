import re

from flask import Flask, jsonify, request
from database import *
import json

TRAVERSAL_REGEX = re.compile(r'(^|[\\/])(\.\.?)([\\/]|$)')

app = Flask(__name__)
@app.route("/upload", methods=["POST"])
def upload():
    file_name = request.form.get("name")
    file_location = request.form.get("location")
    file_type = request.form.get("type")
    if 'file' not in request.files:
        return jsonify(error="No file provided"), 400
    if file_location is None or file_location == "":
        return jsonify(error="No location provided"), 400
    if file_name is None or file_name == "":
        return jsonify(error="No name provided"), 400
    if file_type is None or file_type == "":
        return jsonify(error="No type provided"), 400
    if TRAVERSAL_REGEX.search(file_location) or TRAVERSAL_REGEX.search(file_name):
        return jsonify(error="Traversal is not allowed"), 400

    try:
        new_file = File(name=file_name,location=file_location,type="file",created=datetime.datetime.now(),updated=datetime.datetime.now())
        new_file.save()
    except IntegrityError(err):
        print()
        return jsonify(error="Couldn't create file")
