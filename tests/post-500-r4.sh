#!/bin/bash

DIR="../../synthetic/fhir/r4/500"

for f in $DIR/*.json; do
    curl -X POST http://localhost:3000/ \
        --header "Content-Type: application/json" \
        -d @$f
    echo
done
