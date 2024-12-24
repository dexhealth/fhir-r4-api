#!/bin/bash

FHIR_URL="http://localhost:3000"
RES_URL=$FHIR_URL"/Patient"

### create

curl -s -X POST -H "Content-Type: application/json" -d @patients/david_beckham.json $RES_URL
curl -s -X POST -H "Content-Type: application/json" -d @patients/david_beckham.json $RES_URL

### retrieve

# curl -s http://localhost:3000/Patient
# curl -s http://localhost:3000/Patient?gender=female
# curl -s http://localhost:3000/Patient?family=Doe
# curl -s http://localhost:3000/Patient?postalCode=33139 | jq '.[].name[]'
# curl -s http://localhost:3000/Patient?state=Florida | jq '.[].telecom[]'
# curl -s http://localhost:3000/Patient?birthdate=1998-11-29 | jq '.[].name'

# curl -s http://localhost:3000/Patient?birthdate=gt1987-03-13 | jq '.[].name'

# curl -s http://localhost:3000/Patient?phone=14169001222 | jq '.[].name[]'
# curl -s http://localhost:3000/Patient?email=dbeckham@intermiamicf.com | jq '.[].name'
# curl -s http://localhost:3000/Patient?identifier=98765432 | jq '.[].identifier[]'
# curl -s http://localhost:3000/Patient?given=David,Jane| jq '.[].name[].given'
# curl -s "http://localhost:3000/Patient?country=BRA&given=David,Jane" | jq '.[].name[].given'
# curl -s "http://localhost:3000/Patient?given=David,Jane&country=Brazil" | jq '.[].name[].given'

# curl -s "http://localhost:3000/Patient?identifier=98765432&country=USA" | jq '.[].name'

# curl -s "http://localhost:3000/Patient?gender=female&_sort=name.family" | jq '.[].name[].family'

# curl -s "http://localhost:3000/Patient?phone=1416" | jq '.[].name'
# curl -s "http://localhost:3000/Patient?phone:exact=14169001222" | jq '.[].name'
# curl -s "http://localhost:3000/Patient?phone:contains=234" | jq '.[].name'

# jq '.[] | select(.gender=="female")'
# jq '.[] | select(.name[].family=="Beckham")'

#echo "curl -s "$RES_URL"/"$DAVID_ID" | jq '.'"
