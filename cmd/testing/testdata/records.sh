#!/bin/bash

FHIR_URL="http://localhost:3000"
RES_URL=$FHIR_URL"/"

### create
i=0
for file in 100_Patient_Records/*; do 
    if [ -f "$file" ]; then 
        temp=$(curl -s -X POST -H "Content-Type: application/json" -d @"$file" $RES_URL)
        echo "$temp"
        if [ ${#temp} -ge 100 ]; then 
          echo "$file"
          ((i += 1)) 
        fi
    fi 
done
echo $i

#DAVID_ID=$(curl -s -X POST -H "Content-Type: application/json" -d @patients/david_beckham.json $RES_URL)
