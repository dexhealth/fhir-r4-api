#!/bin/bash

DIR="../../synthetic/fhir/r4/1"

# entire resource
curl -X POST http://localhost:3000/ \
    --header "Content-Type: application/json" \
    -d @"${DIR}/Dorine689_O'Reilly797_a82d1937-b149-0805-f2dd-521add5afa2f.json"

# 1 resource in bundle, the Patient
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_1.json"

# 10 resources in the bundle
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_10.json"

# 9 resources in the bundle, no Encounter
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_NoEncounter.json"

# 20 resources in the bundle
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_20.json"

# 19 resources in the bundle, No encounter
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_19_NoEncounters.json"

# Encounter only
# curl -X POST http://localhost:3000/ \
#     --header "Content-Type: application/json" \
#     -d @"${DIR}/Dorine689_O'Reilly797_EncounterOnly.json"

# List of errors
#
# /Encounter

# List of passess
#
# /Patient
# /Observation
# /Condition
# /Provenance