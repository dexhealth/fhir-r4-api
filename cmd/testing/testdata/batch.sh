#!/bin/bash

FHIR_URL="http://localhost:3000"
RES_URL=$FHIR_URL"/"

### batch

BATCH_GET=$(curl -s -X POST -H "Content-Type: application/json" -d @bundle/get.json $RES_URL)

### Response

echo $BATCH_GET | jq '.'

