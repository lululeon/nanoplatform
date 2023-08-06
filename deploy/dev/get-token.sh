#! /usr/bin/bash

 curl -H 'Content-Type: application/json' \
      -d '{ "formFields": [{ "id": "email", "value": "carol.test@gmail.com" }, { "id": "password", "value": "Testing123" }]}' \
      -X POST \
      https://auth-service:7567/signin
