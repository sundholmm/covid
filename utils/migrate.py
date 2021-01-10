from requests import get, post
from pathlib import Path
from json import load, dumps, dump
from os import remove

# This script downloads and corrects the covid data in JSON format from opendata.ecdc.europa.eu
# and migrates it to the running covid application

records = Path('./data/records.json')
if not records.is_file():
    r = get('https://opendata.ecdc.europa.eu/covid19/casedistribution/json/')
    records.write_bytes(r.content)

with open(records, 'r') as data_file_read:
    data = load(data_file_read)
    for record in data['records']:
        if record['notification_rate_per_100000_population_14-days'] == '':
            record['notification_rate_per_100000_population_14-days'] = None

remove(records)

with open(records, 'w') as data_file_write:
    dump(data, data_file_write, indent=4)
    post('http://localhost:8080/api/v1/records', data = dumps(data))