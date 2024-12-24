#!/bin/bash

FHIR_URL="http://localhost:3000"
RES_URL=$FHIR_URL"/Practitioner"

### create

BRIAN_ID=$(curl -s -X POST -H "Content-Type: application/json" -d @practitioners/brian.json $RES_URL)
FRANK_ID=$(curl -s -X POST -H "Content-Type: application/json" -d @practitioners/frank.json $RES_URL)

### retrieve
# curl http://localhost:3000/Practitioner

curl -s $RES_URL"/"$BRIAN_ID | jq '.name'
curl -s $RES_URL"/"$FRANK_ID | jq '.name'

### search active

curl $RES_URL"?active=true"