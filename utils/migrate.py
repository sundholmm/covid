from requests import get, post
from pathlib import Path
from json import load, dumps

# This script downloads the covid data in JSON format from opendata.ecdc.europa.eu
# and migrates it to the running (!) covid application

records = Path('./data/records.json')
if not records.is_file():
    r = get('https://opendata.ecdc.europa.eu/covid19/casedistribution/json/')
    records.write_bytes(r.content)

with open(records) as data_file:
    data = load(data_file)
    post('http://localhost:8080/api/v1/records', data = dumps(data))