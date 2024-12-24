# Using Azure Health Data Services

To learn FHIR, you can use the Azure FHIR API.  This is for testing & dev only.

- Limit API calls to 2,000 per day.  Do not access the service in an experimental "loop", as a bug can create very high cost charges from Microsoft.
- Microsoft will lock the API if too many requests for generating access tokens are executed (to https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/token).  They don't specify what the limits are, but note that an Access token is valid for 1 hour, so re-use the token until it expires, then generate a new one when needed.  You should be able to generate at least 100 access tokens per day.

## Use

Load the Azure environment variables into your session.

```sh
# from the ./testdata/azure directory ...
source .env
```

To get an Access Token loaded into your sessions environment variable, you can run ...

```sh
export ACCESS_TOKEN=$(curl --silent \
    --request POST https://login.microsoftonline.com/${AZURE_FHIR_TENANT_ID}/oauth2/token \
    --header "Content-Type: application/x-www-form-urlencoded" \
    --data "grant_type=client_credentials" \
    --data "client_id=${AZURE_FHIR_CLIENT_ID}" \
    --data "client_secret=${AZURE_FHIR_CLIENT_SECRET}" \
    --data "resource=${AZURE_FHIR_URL}" \
    | jq --raw-output '.access_token')
```

### Create a patient

Request

```sh
curl -s \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer ${ACCESS_TOKEN}" \
	-d @../patients/david_beckham.json \
	-X POST ${AZURE_FHIR_URL}/Patient | jq '.'
```

Response

```json
{
  "resourceType": "Patient",
  "id": "ab30ee70-fdbf-4a79-b933-98fa7decb48f",
  "meta": {
    "versionId": "1",
    "lastUpdated": "2024-04-19T14:49:30.84+00:00"
  },
  "identifier": [
    {
      "system": "OHIP",
      "value": "123456789"
    }
  ],
  "active": true,
  "name": [
    {
      "family": "Beckham",
      "given": [
        "David",
        "Robert",
        "Joseph"
      ]
    }
  ],
  "telecom": [
    {
      "system": "phone",
      "value": "13058880000",
      "use": "mobile",
      "rank": 1
    },
    {
      "system": "phone",
      "value": "14169001222",
      "use": "home",
      "rank": 2
    },
    {
      "system": "email",
      "value": "dbeckham@intermiamicf.com",
      "use": "work"
    }
  ],
  "gender": "male",
  "birthDate": "1974-12-25",
  "address": [
    {
      "use": "home",
      "type": "both",
      "line": [
        "130 Palm Ave"
      ],
      "city": "Miami Beach",
      "state": "Florida",
      "postalCode": "33139",
      "country": "USA"
    }
  ],
  "generalPractitioner": [
    {
      "type": "Practitioner",
      "identifier": {
        "system": "CPSO",
        "value": "789"
      }
    }
  ]
}
```

### Get a patient

Request

```sh
# replace patient_id with the one created
curl -s ${AZURE_FHIR_URL}/Patient/ab30ee70-fdbf-4a79-b933-98fa7decb48f \
	-H "Authorization: Bearer ${ACCESS_TOKEN}" | jq '.'
```

Response

```json
{
  "resourceType": "Patient",
  "id": "ab30ee70-fdbf-4a79-b933-98fa7decb48f",
  "meta": {
    "versionId": "1",
    "lastUpdated": "2024-04-19T14:49:30.84+00:00"
  },
  "identifier": [
    {
      "system": "OHIP",
      "value": "123456789"
    }
  ],
  "active": true,
  "name": [
    {
      "family": "Beckham",
      "given": [
        "David",
        "Robert",
        "Joseph"
      ]
    }
  ],
  "telecom": [
    {
      "system": "phone",
      "value": "13058880000",
      "use": "mobile",
      "rank": 1
    },
    {
      "system": "phone",
      "value": "14169001222",
      "use": "home",
      "rank": 2
    },
    {
      "system": "email",
      "value": "dbeckham@intermiamicf.com",
      "use": "work"
    }
  ],
  "gender": "male",
  "birthDate": "1974-12-25",
  "address": [
    {
      "use": "home",
      "type": "both",
      "line": [
        "130 Palm Ave"
      ],
      "city": "Miami Beach",
      "state": "Florida",
      "postalCode": "33139",
      "country": "USA"
    }
  ],
  "generalPractitioner": [
    {
      "type": "Practitioner",
      "identifier": {
        "system": "CPSO",
        "value": "789"
      }
    }
  ]
}
```

### Get all patients

Request

```sh
curl -s ${AZURE_FHIR_URL}/Patient -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq '.'
```

Response

```json
{
  "resourceType": "Bundle",
  "id": "f97c5e0753007e169828443fb033e0c2",
  "meta": {
    "lastUpdated": "2024-04-19T14:54:10.0576775+00:00"
  },
  "type": "searchset",
  "link": [
    {
      "relation": "self",
      "url": "https://surgidex-dev.fhir.azurehealthcareapis.com/Patient"
    }
  ],
  "entry": [
    {
      "fullUrl": "https://surgidex-dev.fhir.azurehealthcareapis.com/Patient/ab30ee70-fdbf-4a79-b933-98fa7decb48f",
      "resource": {
        "resourceType": "Patient",
        "id": "ab30ee70-fdbf-4a79-b933-98fa7decb48f",
        "meta": {
          "versionId": "1",
          "lastUpdated": "2024-04-19T14:49:30.84+00:00"
        },
        "identifier": [
          {
            "system": "OHIP",
            "value": "123456789"
          }
        ],
        "active": true,
        "name": [
          {
            "family": "Beckham",
            "given": [
              "David",
              "Robert",
              "Joseph"
            ]
          }
        ],
        "telecom": [
          {
            "system": "phone",
            "value": "13058880000",
            "use": "mobile",
            "rank": 1
          },
          {
            "system": "phone",
            "value": "14169001222",
            "use": "home",
            "rank": 2
          },
          {
            "system": "email",
            "value": "dbeckham@intermiamicf.com",
            "use": "work"
          }
        ],
        "gender": "male",
        "birthDate": "1974-12-25",
        "address": [
          {
            "use": "home",
            "type": "both",
            "line": [
              "130 Palm Ave"
            ],
            "city": "Miami Beach",
            "state": "Florida",
            "postalCode": "33139",
            "country": "USA"
          }
        ],
        "generalPractitioner": [
          {
            "type": "Practitioner",
            "identifier": {
              "system": "CPSO",
              "value": "789"
            }
          }
        ]
      },
      "search": {
        "mode": "match"
      }
    }
  ]
}
```