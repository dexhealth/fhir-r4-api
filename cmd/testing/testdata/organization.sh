#!/bin/bash

FHIR_URL="http://localhost:3000"
RES_URL=$FHIR_URL"/Organization"

### create

hl=$(curl -s -X POST -H "Content-Type: application/json" -d @organizations/hl7.json $RES_URL)

### retrieve

echo "curl -s "$RES_URL"/"$hl" | jq '.'"
