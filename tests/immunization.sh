#!/bin/bash

DIR="../../synthetic/fhir/r4/1"

# entire resource containing an immunization resource
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Abram53_Collier206_d0e71407-a63a-8a4a-a4dd-ee4db4c9f4a5.json"

# Just patient and 1 immunization resource, with occurenceDateTime field
curl -X POST http://localhost:3000/ \
    --header "Content-Type: application/json" \
    -d @"${DIR}/immunization.json"


# https://www.hl7.org/fhir/immunization.html
# there is either an occurenceDateTime or occurenceString, there must be one
# but there isn't an occurence field in immunization