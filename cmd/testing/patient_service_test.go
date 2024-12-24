package patient_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-test/deep"

	"github.com/dexhealth/fhir/pkg/core/search/lookup"
	"github.com/dexhealth/fhir/pkg/core/types/fhirErrors"
	"github.com/dexhealth/fhir/pkg/r5/datastore/mongodb"
	"github.com/dexhealth/fhir/pkg/r5/models"
	"github.com/dexhealth/fhir/pkg/r5/resource/patient"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var svc *patient.Service
var ctrl *patient.Controller

const patientDir = "testdata/patients/"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

/*
	go test -v ./internal/resource/patient '-run=Search'
	go test -v ./internal/resource/patient '-run=Search/Given'
	go test -v ./internal/resource/patient '-run=Search'
	clear; go test -v ./internal/resource/patient '-run=Search'
	clear; go test -v ./internal/resource/patient '-run=Search/Active'
	clear; go test -v ./internal/resource/patient '-run=Search/City'
*/

// setup datastore and seed the test data
func setup() error {
	var err error
	ds, err := mongodb.NewDatastore()
	if err != nil {
		fmt.Printf("could not connect to datastore, %s", err)
		return err
	}

	// init loookup table
	lookupTable, err := lookup.NewSearchTable("../../../internal/gen/definitions/r5/search-parameters.json")
	if err != nil {
		log.Errorf("loading lookup table, %s", err)
		return err
	}
  host := os.Getenv("TEST_HOST")
  port := os.Getenv("TEST_PORT")
  fmt.Printf("Running with host: %s, and port: %s", host, port)
  url  := host+":"+port

	// initialize a new service
	svc = patient.NewService(ds, url, lookupTable)

	// initialize a new controller
	ctrl = patient.NewController(ds, url, lookupTable)

	return nil
}

func teardown() {}

func TestCreateService(t *testing.T) {

	// test cases
	// ------------------------------------------------------
	// simple patient document
	// complex patient document with several nested structs
	// empty document
	// validation check: missing name
	// validation check: incorrect date format for birthDate

	// setup, load test data

	var tt = []struct {
		name      string
		file      string
		id        *string
		wantError error
	}{
		{"Simple Patient", "jane_doe.json", nil, nil},
		{"Complex Patient", "david_beckham.json", nil, nil},
		// {"Missing Name", "no_name.json", nil, fhirErrors.ErrNotValid},
		// {"Empty Patient", "empty.json", nil, fhirErrors.ErrNotValid},
	}

	// run tests, update id field for later clean up

	for i, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			jsonFile, err := os.Open(patientDir + tc.file)
			if err != nil {
				t.Errorf("could not open test case file %s, %s", tc.file, err)
			}
			defer jsonFile.Close()

			byteValue, err := io.ReadAll(jsonFile)
			if err != nil {
				t.Errorf("could not read test case file %s, %s", tc.file, err)
			}

			var patient models.Patient
			json.Unmarshal(byteValue, &patient)

			tt[i].id, err = svc.Create(patient)
			if !errors.Is(err, tc.wantError) {
				t.Errorf("want error: %v, have %v", tc.wantError, err)
			}

			// if expecting success, ensure an id was returned to confirm operation
			if tc.wantError == nil {
				id := tt[i].id
				if len(*id) > 0 {

					// confirm that it exists in the datastore
					_, err := svc.Retrieve(id)
					if err != nil {
						t.Errorf("could not retrieve patient using id %s, %s", *id, err)
					}

					// show me
					// name, _ := json.MarshalIndent(p.Name, "", "  ")
					// t.Logf("%s", name)
				} else {
					t.Errorf("expected id, got empty string")
				}
			}
		})
	}

	// cleanup, remove inserted data

	for _, tc := range tt {
		if tc.id != nil {
			err := svc.Delete(tc.id)
			if err != nil {
				t.Errorf("error during cleanup, %s", err)
			}
		}
	}
}

func TestDeleteService(t *testing.T) {

	// load test data

	invalidId := "this-is-invalid-uuid-format"
	fakeId := primitive.NewObjectID().Hex()
	emptyId := ""

	var table = []struct {
		name string
		file string
		want error
		id   *string
	}{
		{"Valid Patient", "david_beckham.json", nil, nil},
		{"Invalid ID", "", primitive.ErrInvalidHex, &invalidId},
		{"ID Not Found", "", fhirErrors.ErrNotFound, &fakeId},
		{"Missing ID", "", primitive.ErrInvalidHex, &emptyId},
	}

	// insert patients (if they have a json file), update id field

	for i, tc := range table {
		if len(tc.file) > 0 {
			jsonFile, err := os.Open(patientDir + tc.file)
			if err != nil {
				t.Errorf("could not open test case file %s, %s", tc.file, err)
			}
			defer jsonFile.Close()

			byteValue, err := io.ReadAll(jsonFile)
			if err != nil {
				t.Errorf("could not read test case file %s, %s", tc.file, err)
			}

			var patient models.Patient
			json.Unmarshal(byteValue, &patient)

			table[i].id, err = svc.Create(patient)
			if err != tc.want {
				t.Errorf("error creating patient, %s", err)
			}
		}
	}

	// run tests

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.Delete(tc.id)
			if !errors.Is(err, tc.want) {
				t.Errorf("want: %+v, have: %+v", tc.want, err)
			}
			// check that the patient was actually deleted from the datastore
			if errors.Is(nil, tc.want) {
				p, err := svc.Retrieve(tc.id)
				if !errors.Is(err, fhirErrors.ErrNotFound) {
					if p != nil {
						t.Errorf("patient not deleted, %s", *p.Id)
					}
				}
			}
		})
	}

	// clean up, remove sample inserts

}

func TestUpdateService(t *testing.T) {

	// test cases, tc
	// -------------
	// success_replace
	// invalid_id
	// id_not_found
	// missing_id
	// id_mismatch

	// load test cases (tc) in test table (tt)

	var tt = []struct {
		name                string          /* name of test */
		originalPatientFile string          /* patient before update*/
		id                  *string         /* patient id to update */
		wantPatientFile     string          /* patient after update */
		wantPatient         *models.Patient /* expected patient after update */
		wantError           error           /* anticipated / expected error */
	}{
		{"success replace", "david_beckham.json", nil, "victoria_beckham.json", nil, nil},
	}

	// insert patients (if thye have a json file), update id and wantPatient

	for i, tc := range tt {
		// only update id and wantPatient if original patient file provided
		if len(tc.originalPatientFile) > 0 && len(tc.wantPatientFile) > 0 {
			// load original patient
			jsonFile, err := os.Open(patientDir + tc.originalPatientFile)
			if err != nil {
				t.Errorf("could not open test case file %s, %s", tc.originalPatientFile, err)
			}
			defer jsonFile.Close()

			byteValue, err := io.ReadAll(jsonFile)
			if err != nil {
				t.Errorf("could not read test case file %s, %s", tc.originalPatientFile, err)
			}

			var patient models.Patient
			json.Unmarshal(byteValue, &patient)

			// will require id for cleanup of datastore when tests complete
			tt[i].id, err = svc.Create(patient)
			if err != tc.wantError {
				t.Errorf("error creating patient, %s", err)
			}

			// load expected patient, after update
			jsonFileExpected, err := os.Open(patientDir + tc.wantPatientFile)
			if err != nil {
				t.Errorf("could not open test case file %s, %s", tc.wantPatientFile, err)
			}
			defer jsonFileExpected.Close()

			byteValueExpected, err := io.ReadAll(jsonFileExpected)
			if err != nil {
				t.Errorf("could not read test case file %s, %s", tc.wantPatientFile, err)
			}

			var patientExpected models.Patient
			json.Unmarshal(byteValueExpected, &patientExpected)

			tt[i].wantPatient = &patientExpected
		}
	}

	// run tests

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.Update(tc.id, tc.wantPatient)
			if !errors.Is(err, tc.wantError) {
				t.Errorf("want error: %+v, have: %+v", tc.wantError, err)

				expectedName, _ := json.MarshalIndent(tc.wantPatient.Name, "", "  ")
				t.Errorf("\n%s", expectedName)
			}

			// check that the patient was actually updated in the datastore
			if errors.Is(nil, tc.wantError) {
				havePatient, err := svc.Retrieve(tc.id)
				if err != nil {
					t.Errorf("error verifying patient update, %s", err)
				}
				// don't compare id's, remove from retrieved patient
				havePatient.Id = nil
				havePatient.Latest = nil
				//lat := false
				ver := 2
				tc.wantPatient.VersionHistory = &ver
				diff := deep.Equal(havePatient, tc.wantPatient)
				if diff != nil {
					t.Errorf("want patient: %+v, \nhave: %+v\n", tc.wantPatient, havePatient)
					t.Errorf("DIFF: %+v\n", diff)

					// show me
					haveName, _ := json.MarshalIndent(havePatient.Name, "", "  ")
					t.Errorf("\n%s", haveName)
					t.Errorf("\nid used is %s", *tc.id)
				}
			}
		})
	}

	// cleanup, remove sample patients

	for _, tc := range tt {
		if tc.id != nil {
			err := svc.Delete(tc.id)
			if err != nil {
				t.Errorf("error during cleanup, %s", err)
			}
		}
	}

}
