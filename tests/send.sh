#! /bin/bash

curl -vvv http://0.0.0.0:5000/send -H "Content-Type: application/json" -d @tests/data.json
