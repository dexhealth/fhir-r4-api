package patient_test

import (
	"encoding/json"

	//"github.com/go-test/deep"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/dexhealth/fhir/pkg/core/types/dto"
	"github.com/dexhealth/fhir/pkg/r5/models"
)

// index of test patients
const (
	TP_EDSON_DO_NASCIMENTO string = "Edson do Nascimento"
	TP_DAVID_BECKHAM       string = "David Beckham"
	TP_VICTORIA_BECKHAM    string = "Victoria Beckham"
	TP_VICTOR_OSIMHEN      string = "Victor Osimhen"
	TP_DAVID_GARCIA        string = "David Garcia"
	TP_JANE_DOE            string = "Jane Doe"
	TP_ANDREAS_BECK        string = "Andreas Beck"
	TP_ROBERT_DEANGELIS    string = "Robert DeAngelis"
)

//var ctrl *patient.Controller

// Contain sample patient data that must be inserted and deleted after each test func
type PatientData struct {
	Id   *string
	File string
	Data *models.Patient
}

func getFilenameFromName(name string) (filename string) {
	filename = strings.Replace(name, " ", "_", -1) + ".json"
	return strings.ToLower(filename)
}

func insertPatientData() (map[string]*PatientData, error) {
	patientMap := make(map[string]*PatientData)

	patientMap[TP_EDSON_DO_NASCIMENTO] = &PatientData{nil, getFilenameFromName(TP_EDSON_DO_NASCIMENTO), nil}
	patientMap[TP_DAVID_BECKHAM] = &PatientData{nil, getFilenameFromName(TP_DAVID_BECKHAM), nil}
	patientMap[TP_VICTORIA_BECKHAM] = &PatientData{nil, getFilenameFromName(TP_VICTORIA_BECKHAM), nil}
	patientMap[TP_VICTOR_OSIMHEN] = &PatientData{nil, getFilenameFromName(TP_VICTOR_OSIMHEN), nil}
	patientMap[TP_DAVID_GARCIA] = &PatientData{nil, getFilenameFromName(TP_DAVID_GARCIA), nil}
	patientMap[TP_JANE_DOE] = &PatientData{nil, getFilenameFromName(TP_JANE_DOE), nil}
	patientMap[TP_ANDREAS_BECK] = &PatientData{nil, getFilenameFromName(TP_ANDREAS_BECK), nil}

	for _, p := range patientMap {
		jsonFile, err := os.Open(patientDir + p.File)
		if err != nil {
		}
		defer jsonFile.Close()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
		}

		var patient models.Patient
		json.Unmarshal(byteValue, &patient)

		p.Data = &patient

		p.Id, err = svc.Create(*p.Data)
		if err != nil {
		}
	}

	return patientMap, nil
}

func deletePatientData(patientMap map[string]*PatientData) error {
	// fmt.Printf("found %d patients to delete from database\n", len(patientMap))
	for _, p := range patientMap {
		// fmt.Printf("deleting %s\n", *p.Id)
		err := svc.Delete(p.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestPagingController(t *testing.T) {
	// [US-1.9: Patient Search Results Paging](https://github.com/dexhealth/fhir/issues/10)
	// go test -v ./internal/resource/patient '-run=Paging'

	/******************************  N O T E S  ******************************

	_count is defined as an instruction to the server regarding how many resources should be returned in a single page
	_count is the pageSize
	ref: https://cloud.google.com/healthcare-api/docs/how-tos/fhir-search#pagination_and_sorting

	docker exec -it fhir-mongo-1 mongosh --username root

	curl -s 'http://localhost:3000/Patient?_sort=family,given&_count=2'
	curl -s 'http://localhost:3000/Patient?_sort=family,given&_count=2&page=1'
	curl -s 'http://localhost:3000/Patient?_sort=family,given&_count=2&page=2'

	jq '.entry[].name'
	jq '.links[]'

	- result always comes up with &page=1, appears the pagination logic is missing

	*************************************************************************/

	// insert sample data

	patientMap, err := insertPatientData()
	if err != nil {
		t.Error(err)
	}

	// define tests, validate response bundle containing patients and links

	var tt = []struct {
		name                   string
		whenQuery              string
		expectStatusCode       int
		expectedResponseBundle dto.Bundle
	}{
		{
			"US-1.9 Page Size 2 Default", "?_sort=family,given&_count=2", http.StatusOK,
			dto.Bundle{
				Entry: []models.Patient{
					*patientMap[TP_ANDREAS_BECK].Data,
					*patientMap[TP_VICTORIA_BECKHAM].Data,
				},
				Links: []dto.Link{
					{Relation: "first", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "last", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "next", Url: "_count=2&_sort=family,given&page=2"},
					{Relation: "self", Url: "_count=2&_sort=family,given&page=1"},
				},
			},
		}, {
			"US-1.9 Page Size 2 Page 1", "?_sort=family,given&_count=2&page=1", http.StatusOK,
			dto.Bundle{
				Entry: []models.Patient{
					*patientMap[TP_ANDREAS_BECK].Data,
					*patientMap[TP_VICTORIA_BECKHAM].Data,
				},
				Links: []dto.Link{
					{Relation: "first", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "last", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "next", Url: "_count=2&_sort=family,given&page=2"},
					{Relation: "self", Url: "_count=2&_sort=family,given&page=1"},
				},
			},
		}, {
			"US-1.9 Page Size 2 Page 2", "?_sort=family,given&_count=2&page=2", http.StatusOK,
			dto.Bundle{
				Entry: []models.Patient{
					*patientMap[TP_DAVID_BECKHAM].Data,
					*patientMap[TP_JANE_DOE].Data,
				},
				Links: []dto.Link{
					{Relation: "first", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "last", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "previous", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "next", Url: "_count=2&_sort=family,given&page=3"},
					{Relation: "self", Url: "_count=2&_sort=family,given&page=2"},
				},
			},
		}, {
			"US-1.9 Page Size 2 Page 3", "?_sort=family,given&_count=2&page=3", http.StatusOK,
			dto.Bundle{
				Entry: []models.Patient{
					*patientMap[TP_DAVID_GARCIA].Data,
					*patientMap[TP_VICTOR_OSIMHEN].Data,
				},
				Links: []dto.Link{
					{Relation: "first", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "last", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "previous", Url: "_count=2&_sort=family,given&page=2"},
					{Relation: "next", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "self", Url: "_count=2&_sort=family,given&page=3"},
				},
			},
		}, {
			"US-1.9 Page Size 2 Page 4", "?_sort=family,given&_count=2&page=4", http.StatusOK,
			dto.Bundle{
				Entry: []models.Patient{
					*patientMap[TP_EDSON_DO_NASCIMENTO].Data,
				},
				Links: []dto.Link{
					{Relation: "first", Url: "_count=2&_sort=family,given&page=1"},
					{Relation: "last", Url: "_count=2&_sort=family,given&page=4"},
					{Relation: "previous", Url: "_count=2&_sort=family,given&page=3"},
					{Relation: "self", Url: "_count=2&_sort=family,given&page=4"},
				},
			},
		},
	}

	// run tests

	e := echo.New()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/Patient"+tc.whenQuery, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ctrl.Search(c)
			if err != nil {
				t.Errorf("error searching test case %s, %s", tc.whenQuery, err)
			}

			assert.Equal(t, tc.expectStatusCode, rec.Code)

			// expect patient(s) as json string in response body
			if tc.expectStatusCode <= 300 {
				var responseBundle dto.Response
				err = json.Unmarshal(rec.Body.Bytes(), &responseBundle)
				if err != nil {
					t.Errorf("error converting response body to response bundle, %s", err)
				}

				// the number of links returned must equal the number expected
				if len(tc.expectedResponseBundle.Links) != len(responseBundle.Links) {
					assert.Equal(t, len(tc.expectedResponseBundle.Links), len(responseBundle.Links))
				} else {
					for _, expectedLinks := range tc.expectedResponseBundle.Links {
						if !containsLink(tc.expectedResponseBundle.Links, expectedLinks) {
							t.Errorf("missing data in response body")
						}
					}
				}

				expectedPatients := tc.expectedResponseBundle.Entry.([]models.Patient)

				// convert response bundle entry to array of Patients
				var havePatients []models.Patient
				for _, p := range responseBundle.Entry {
					// convert map to json string
					jsonStr, err := json.Marshal(p.Resource)
					if err != nil {
						t.Error(err)
					}
					// convert json string to struct
					var havePatient models.Patient
					if err := json.Unmarshal(jsonStr, &havePatient); err != nil {
						t.Error(err)
					}
					havePatient.Id = nil
					havePatient.VersionHistory = nil
					havePatient.Latest = nil
					havePatient.Oid = nil
					havePatients = append(havePatients, havePatient)
				}

				// the number of patients returned must equal the number expected
				if len(expectedPatients) != len(havePatients) {
					assert.Equal(t, len(expectedPatients), len(havePatients))
				} else {
					// compare/match patients, in expected order
					// removing assert.Equal, because the debug printf's are too messy for human reading
					// assert.Equal(t, expectedPatients, havePatients)
					if diff := deep.Equal(expectedPatients, havePatients); diff != nil {
						t.Error(diff)
						// t.Error("nto equql")
					}

				}

			}
		})
	}

	// remove sample data

	err = deletePatientData(patientMap)
	if err != nil {
		t.Error(err)
	}

}

// go test -v ./internal/resource/patient '-run=Sort'
// go test -v ./internal/resource/patient '-run=Sort/Identifier'
func TestSortController(t *testing.T) {

	// insert sample data

	patientMap, err := insertPatientData()
	if err != nil {
		t.Error(err)
	}

	// load tests

	var tt = []struct {
		name             string
		whenQuery        string
		expectStatusCode int
		expectPatients   []string
	}{
		{"US-1.7 Search Identifier Sort Identifier Ascending", "?identifier=123456789&_sort=identifier", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM}},
		{"US-1.7 Search Identifier Sort Family Ascending", "?identifier=123456789&_sort=family", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
		{"US-1.7 Search Identifier Sort Family Decending", "?identifier=123456789&_sort=-family", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM}},

		{"US-1.7 Search Identifier Code Sort Identifier Ascending", "?identifier=OHIP|123456789&_sort=identifier", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"US-1.7 Search Identifier Code Sort Family Ascending", "?identifier:contains=OHIP|456&_sort=family", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_DAVID_GARCIA}},
		{"US-1.7 Search Identifier Code Sort Family Ascending", "?identifier:contains=OHIP|456&_sort=-family", http.StatusOK, []string{TP_DAVID_GARCIA, TP_DAVID_BECKHAM}},

		// ascending birth date implies the sorting will be from older to younger
		// descending birth data implies the sorting will from from youger to older
		{"US-1.7 Search Family Sort Birth Date Ascending", "?family=beckham&_sort=birthdate", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"US-1.7 Search Family Sort Birth Date Descending", "?family=beckham&_sort=-birthdate", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},

		// ascending sorts patients by their "lowest" Given name {Andreas, Caroline, David}
		// descending sorts patients by their "highest" Given name {Victoria, Robert, Andreas}
		{"US-1.7 Search Family Sort Given Ascending", "?family=beck&_sort=given", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"US-1.7 Search Family Sort Given Descending", "?family=beck&_sort=-given", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM, TP_ANDREAS_BECK}},

		{"US-1.7 Search Country Sort Birth Date Ascending", "?address-country=USA&_sort=birthdate", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"US-1.7 Search Country Sort Birth Date Descending", "?address-country=USA&_sort=-birthdate", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
		{"US-1.7 Search Active Sort Identifier Acending", "?active=true&_sort=identifier", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
		{"US-1.7 Search Active Sort Identifier Decending", "?active=true&_sort=-identifier", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM, TP_EDSON_DO_NASCIMENTO}},

		{"US-1.7 Search Active Family Identifier Sort Given Ascending", "?active=true&identifier=456&_sort=family", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
		{"US-1.7 Search Active Family Identifier Sort Given Decending", "?active=true&identifier=456&_sort=-family", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
	}

	// run tests

	e := echo.New()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/Patient"+tc.whenQuery, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ctrl.Search(c)
			if err != nil {
				t.Errorf("error searching test case %s, %s", tc.whenQuery, err)
			}

			assert.Equal(t, tc.expectStatusCode, rec.Code)

			// expect patient(s) as json string in response body
			if tc.expectStatusCode <= 300 {
				///////
				var responseBundle dto.Response
				err = json.Unmarshal(rec.Body.Bytes(), &responseBundle)
				if err != nil {
					t.Errorf("error converting response body to slice of structs, %s", err)
				}
				var responsePatients []*models.Patient
				for _, patient := range responseBundle.Entry {
					jsonBytes, err := json.Marshal(patient.Resource)
					if err != nil {
						t.Errorf("Error marshalling JSON: %s", err)
					}

					// Unmarshal JSON into struct
					var p models.Patient
					err = json.Unmarshal(jsonBytes, &p)
					if err != nil {
						t.Errorf("Error unmarshalling JSON: %s", err)
					}
					p.Id = nil
					p.VersionHistory = nil
					p.Latest = nil
					p.Oid = nil
					responsePatients = append(responsePatients, &p)
				}
				// the number of patients returned must equal the number expected
				if len(tc.expectPatients) != len(responsePatients) {
					assert.Equal(t, len(tc.expectPatients), len(responsePatients))
				} else {
					// compare/match patients, in expected order
					for i, expectedPatient := range tc.expectPatients {
						assert.Equal(t, patientMap[expectedPatient].Data, responsePatients[i])
					}
				}

			}

		})
	}

	// remove sample data

	err = deletePatientData(patientMap)
	if err != nil {
		t.Error(err)
	}
}

// func TestPagingController(t *testing.T){}

// go test -v ./internal/resource/patient '-run=Search' | grep FAIL:
// go test ./internal/resource/patient '-run=Search'
//
// go test -v ./internal/resource/patient '-run=Search/Given'
// go test -v ./internal/resource/patient '-run=Search/Given_Single_Result'
// go test -v ./internal/resource/patient '-run=Search/Bad_Request'
func TestSearchController(t *testing.T) {

	// load sample patients into database and record id for later cleanup

	patientMap, err := insertPatientData()
	if err != nil {
		t.Error(err)
	}

	// load test data

	var tt = []struct {
		name             string
		whenQuery        string
		expectStatusCode int
		expectPatients   []string
	}{

		// [US-1.3: Patient Search Exact w/Multivalues](https://github.com/dexhealth/fhir/issues/4)
		// go test -v ./internal/resource/patient '-run=Search/US-1.3'

		{"US-1.3 Given Missing Value", "?given=", http.StatusBadRequest, []string{}},
		{"US-1.3 Given Starts With", "?given=Victor", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_VICTOR_OSIMHEN}},
		{"US-1.3 Given Middle Starts With", "?given=Carol", http.StatusOK, []string{TP_VICTORIA_BECKHAM}},
		{"US-1.3 Given Starts With Case-Insensitive", "?given=victor", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_VICTOR_OSIMHEN}},
		{"US-1.3 Family Starts With", "?family=Beck", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM, TP_ANDREAS_BECK}},
		{"US-1.3 Family Starts With Given Start With", "?family=Beck&given=Vic", http.StatusOK, []string{TP_VICTORIA_BECKHAM}},
		{"US-1.3 Country Family Starts With", "?address-country=DEU&family=Beck", http.StatusOK, []string{TP_ANDREAS_BECK}},
		// INFO: https://simplifier.net/guide/OntarioContextManagement/ConformanceRules?version=current
		//
		//       For this specification, the telecom format is "+[country code]-[area code]-[first three
		//       digits of the phone number]-[remaining digits of the phone number] ext."[extension]".
		//       Note that if country code is not used, the "+" prefix is skipped.
		//
		//       Example: "+1-555-555-5555 ext. 123456"
		//
		// INFO: https://hl7.org/fhir/datatypes.html#ContactPoint
		//
		//       If capturing a phone, fax or similar contact point, the value should be a properly formatted
		//       telephone number according to ITU-T E.123 (Notation for national and international telephone numbers)
		//
		//       https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-E.123-200102-I!!PDF-E&type=items
		//
		// To urlencode +, use %2B
		// [HTML URL Encoding Reference](https://www.w3schools.com/tags/ref_urlencode.asp)
		//
		{"US-1.3 Phone International", "?phone=%2B49071199938420", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.3 Phone National", "?phone=13058880000", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"US-1.3 Email", "?email=andy-beck@@vfb.de", http.StatusOK, []string{TP_ANDREAS_BECK}},

		// [US-1.4: Patient Search Fuzzy w/Union OR](https://github.com/dexhealth/fhir/issues/5)
		// Starts With, Case-insensitive, supports union search filter (i.e. OR)
		// go test -v ./internal/resource/patient '-run=Search/US-1.4'

		{"US-1.4 Given Starts With Lowercase", "?given=andre", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.4 Given Starts With Uppercase", "?given=AND", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.4 Given Starts With Mixed-case", "?given=aNdR", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.4 Given Union Lowercase", "?given=andr,victor", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},
		{"US-1.4 Given Union Uppercase", "?given=ANDRE,VICT", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},
		{"US-1.4 Given Union Mixed-case", "?given=aND,vIcTo", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},
		{"US-1.4 Family Starts With Mixed-case", "?family=bEcK", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		// {"US-1.4 Email Start With", "?email=andy-beck", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.4 Identifier Starts With", "?identifier=123", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
		{"US-1.4 Identifier System Starts With", "?identifier=OHIP|123", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_BECKHAM}},
		// {"US-1.4 Identifier Union", "?identifier=900,837", http.StatusOK, []string{TP_ANDREAS_BECK, TP_DAVID_GARCIA}},
		// {"US-1.4 Identifier System Union", "?identifier=OHIP|123,112233", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_BECKHAM, TP_EDSON_DO_NASCIMENTO}},

		// [US-1.5: Patient Search Fuzzy Match](https://github.com/dexhealth/fhir/issues/6)
		// String search modifiers: Exact, Contains
		// String search, no modifier (default behaviour)): Starts-With
		// go test -v ./internal/resource/patient '-run=Search/US-1.5'

		{"US-1.5 Identifier Contains", "?identifier:contains=89", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_VICTOR_OSIMHEN, TP_DAVID_BECKHAM}},
		{"US-1.5 Identifier Exact", "?identifier:exact=123000789", http.StatusOK, []string{TP_VICTOR_OSIMHEN}},
		// {"US-1.5 Identifier Contains Union", "?identifier:contains=678,009", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_ANDREAS_BECK}},
		// {"US-1.5 Identifier Exact Union", "?identifier:exact=123000789,1244566", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},
		{"US-1.5 Identifier System Contains", "?identifier:contains=OHIP|456", http.StatusOK, []string{TP_DAVID_GARCIA, TP_DAVID_BECKHAM}},
		// {"US-1.5 Identifier System Contains Union", "?identifier:contains=OHIP|1122,0900", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_ANDREAS_BECK}},
		{"US-1.5 Identifier System Exact", "?identifier:exact=OHIP|900900900", http.StatusOK, []string{TP_ANDREAS_BECK}},
		// {"US-1.5 Identifier System Exact Union", "?identifier:exact=OHIP|123000789,837649456", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_GARCIA}},
		{"US-1.5 Family Contains", "?family:contains=ham", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
		{"US-1.5 Given Contains Union", "?given:contains=icto,ante", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
		{"US-1.5 Given Exact", "?given:exact=Victor", http.StatusOK, []string{TP_VICTOR_OSIMHEN}},
		{"US-1.5 Given Exact Union", "?given:exact=Victor,Joseph", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_BECKHAM}},
		{"US-1.5 Given Contains Family Contains", "?given:contains=avi&family:contains=eck", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"US-1.5 Given Exact Family Exact", "?given:exact=David&family:exact=Beckham", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"US-1.5 Family Exact", "?family:exact=Beck", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.5 Family Exact Given Starts With", "?family:exact=Beckham&given=Vic", http.StatusOK, []string{TP_VICTORIA_BECKHAM}},

		// [US-1.6: Patient Search Prefixes](https://github.com/dexhealth/fhir/issues/7)
		// go test -v ./internal/resource/patient '-run=Search/US-1.6'

		{"US-1.6 Birth Date Missing", "?birthdate=", http.StatusBadRequest, []string{}},
		{"US-1.6 Birth Date No Results", "?birthdate=1900-01-01", http.StatusOK, []string{}},
		{"US-1.6 Birth Date No Results Bad Prefix", "?birthdate=pp1987-03-13", http.StatusOK, []string{}},
		{"US-1.6 Birth Date Unsupported Date Format", "?birthdate=eqABCD", http.StatusBadRequest, []string{}},
		{"US-1.6 Birth Date", "?birthdate=1987-03-13", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.6 Birth Date EQ", "?birthdate=eq1987-03-13", http.StatusOK, []string{TP_ANDREAS_BECK}},
		{"US-1.6 Birth Date NE", "?birthdate=ne1987-03-13", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM, TP_DAVID_GARCIA, TP_JANE_DOE}},
		{"US-1.6 Birth Date GT", "?birthdate=gt1987-03-13", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_GARCIA}},
		{"US-1.6 Birth Date GE", "?birthdate=ge1987-03-13", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_DAVID_GARCIA}},
		{"US-1.6 Birth Date LT", "?birthdate=lt1987-03-13", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_JANE_DOE, TP_VICTORIA_BECKHAM}},
		{"US-1.6 Birth Date LE", "?birthdate=le1987-03-13", http.StatusOK, []string{TP_ANDREAS_BECK, TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_JANE_DOE, TP_VICTORIA_BECKHAM}},

		// sa - the resource value Starts After the parameter value
		//    - for date comparison, this seems equivalent to 'gt'
		{"US-1.6 Birth Date SA", "?birthdate=sa1987-03-13", http.StatusOK, []string{TP_VICTOR_OSIMHEN, TP_DAVID_GARCIA}},

		// eb - the resource value Ends Before the parameter value
		//    - for date comparison, this seems equivalent to 'lt'
		{"US-1.6 Birth Date EB", "?birthdate=eb1987-03-13", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO, TP_DAVID_BECKHAM, TP_JANE_DOE, TP_VICTORIA_BECKHAM}},

		{"US-1.6 Birth Date Range", "?birthdate=gt1987-03-13,lt1998-11-29", http.StatusOK, []string{TP_DAVID_GARCIA}},
		{"US-1.6 Birth Date Range Reverse", "?birthdate=lt1998-11-29,gt1987-03-13", http.StatusOK, []string{TP_DAVID_GARCIA}},
		{"US-1.6 Birth Date Range Inclusive", "?birthdate=ge1987-03-13,le1998-11-29", http.StatusOK, []string{TP_ANDREAS_BECK, TP_DAVID_GARCIA, TP_VICTOR_OSIMHEN}},
		{"US-1.6 Birth Date Range Inclusive Reverse", "?birthdate=le1998-11-29,ge1987-03-13", http.StatusOK, []string{TP_ANDREAS_BECK, TP_DAVID_GARCIA, TP_VICTOR_OSIMHEN}},

		// [US-1.8: Patient Search Results Limit](https://github.com/dexhealth/fhir/issues/9)

		{"US-1.8 Search Given Union Limit Missing Value", "?given=andr,victor&_maxresults=", http.StatusBadRequest, []string{}},

		// this it not a case that would be used in a real world query
		// {"US-1.8 Search Given Union Limit Zero", "?given=andr,victor&_maxresults=0", http.StatusOK, []string{}},

		{"US-1.8 Search Given Union Limit Large", "?given=andr,victor&_maxresults=10", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},
		{"US-1.8 Search Given Union Sort Given Limit Same", "?given=andr,victor&_sort=given&_maxresults=3", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTOR_OSIMHEN, TP_VICTORIA_BECKHAM}},

		// the search without _maxresults returns 3 items, but here we want to limit the results to 2
		// go test -v ./internal/resource/patient '-run=Search/US-1.8_Search_Given_Union_Sort_Given_Limit_Small'
		//
		// docker exec -it fhir-mongo-1 mongosh --username root
		// db.patients.find( { "name.family": "Beckham" } ).limit(1)
		// db.patients.deleteMany( {} )
		//
		// curl http://localhost:3000/Patient?given=andr,victor&_sort=given&_maxresults=2
		//
		{"US-1.8 Search Given Union Sort Given Limit Small", "?given=andr,victor&_sort=given&_maxresults=2", http.StatusOK, []string{TP_ANDREAS_BECK, TP_VICTORIA_BECKHAM}},

		// Random test cases

		{"Identifier Single Result", "?identifier=837649456", http.StatusOK, []string{TP_DAVID_GARCIA}},
		{"Identifier Specify System", "?identifier=OHIP|837649456", http.StatusOK, []string{TP_DAVID_GARCIA}},
		{"Given Single Result", "?given=Vic", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_VICTOR_OSIMHEN}},
		{"Given Multivalue OR", "?given=Vic,Jan", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_VICTOR_OSIMHEN, TP_JANE_DOE}},
		{"Given Multivalue Exact OR No Results", "?given:exact=Vic,Jan", http.StatusOK, []string{}},
		{"Given No Results", "?given=Persimony", http.StatusOK, []string{}},
		{"Given Single Result Middle", "?given=Robert", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"Given Multiple Results", "?given=David", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_DAVID_GARCIA}},
		{"Family Multiple Results", "?family=Beckham", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"Given & Family Single Result", "?given=Victoria&family=Beckham", http.StatusOK, []string{TP_VICTORIA_BECKHAM}},
		{"Given & Family No Results", "?given=Victoria&family=Wonder", http.StatusOK, []string{}},
		{"Gender Males", "?gender=male", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_DAVID_GARCIA, TP_EDSON_DO_NASCIMENTO, TP_VICTOR_OSIMHEN, TP_ANDREAS_BECK}},
		{"Gender Bad Request", "?gender=donkey", http.StatusBadRequest, []string{}},
		{"Active Multiple", "?active=true", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
		{"Active Bad Request", "?active=blue", http.StatusBadRequest, []string{}},
		{"Phone Single Result", "?phone=13058880000", http.StatusOK, []string{TP_DAVID_BECKHAM}},
		{"Phone Multiple Results", "?phone=14169001222", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
		{"Phone No Results", "?phone=999", http.StatusOK, []string{}},
		{"Telecom Single Result", "?telecom=dgarcia@fiu.edu", http.StatusOK, []string{TP_DAVID_GARCIA}},
		{"Given Exact", "?given:exact=Victoria", http.StatusOK, []string{TP_VICTORIA_BECKHAM}},
		{"Family Contains", "?family:contains=ham", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_VICTORIA_BECKHAM}},
		{"Birth Date Range", "?birthdate=ge1940-10-29,lt1940-11-01", http.StatusOK, []string{TP_EDSON_DO_NASCIMENTO}},
		{"City", "?address-city=Miami%20Beach", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"State", "?address-state=Florida", http.StatusOK, []string{TP_VICTORIA_BECKHAM, TP_DAVID_BECKHAM}},
		{"Identifier Any", "?identifier=123456789", http.StatusOK, []string{TP_DAVID_BECKHAM, TP_EDSON_DO_NASCIMENTO}},
	}

	// run test cases

	e := echo.New()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// fmt.Printf("preping pateint router test %s\n", tc.name)

			req := httptest.NewRequest(http.MethodGet, "/Patient"+tc.whenQuery, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ctrl.Search(c)
			if err != nil {
				t.Errorf("error searching test case %s, %s", tc.whenQuery, err)
			}

			assert.Equal(t, tc.expectStatusCode, rec.Code)

			// convert body *bytes.Buffer into slice of patients to compare
			if tc.expectStatusCode <= 300 {
				var responseBundle dto.Response
				err = json.Unmarshal(rec.Body.Bytes(), &responseBundle)
				if err != nil {
					t.Errorf("error converting response body to slice of structs, %s", err)
				}
				var responsePatients []models.Patient
				for _, patient := range responseBundle.Entry {
					jsonBytes, err := json.Marshal(patient.Resource)
					if err != nil {
						t.Errorf("Error marshalling JSON: %s", err)
					}

					// Unmarshal JSON into struct
					var p models.Patient
					err = json.Unmarshal(jsonBytes, &p)
					if err != nil {
						t.Errorf("Error unmarshalling JSON: %s", err)
					}
					responsePatients = append(responsePatients, p)
				}

				if len(tc.expectPatients) != len(responsePatients) {
					//for _, p := range responsePatients {
					//  b, _ := json.MarshalIndent(p, "", "  ")
					//  fmt.Printf("%s", b)
					//}
					assert.Equal(t, len(tc.expectPatients), len(responsePatients))
					//t.Errorf("expected %d patients, have %d\n", len(tc.expectPatients), len(responsePatients))

				}

				for _, ep := range tc.expectPatients {
					if !containsNoRef(responsePatients, patientMap[ep].Data) {
						//for _, p := range responsePatients {
						//  b, _ := json.MarshalIndent(p, "", "  ")
						//  fmt.Printf("%s", b)
						//}
						t.Errorf("missing data in response body")
					}
				}
			}

		})
	}

	// remove test patients from database

	err = deletePatientData(patientMap)
	if err != nil {
		//t.Error(err)
		t.Errorf("error deleting sample data, %s", err)
	}
}

// returns true if x is found in a
// func contains(a []*models.Patient, x *models.Patient) bool {
// 	v := 0
// 	for _, n := range a {
// 		n.Id = nil
// 		n.VersionHistory = &v
// 		lat := false
// 		n.Latest = &lat
// 		n.Oid = nil
// 		if cmp.Equal(&x, &n) {
// 			return true
// 		}
// 	}
// 	return false
// }

func containsNoRef(a []models.Patient, x *models.Patient) bool {
	for _, n := range a {
		n.Id = x.Id
		n.VersionHistory = x.VersionHistory
		n.Latest = x.Latest
		n.Oid = x.Oid

		//diff := deep.Equal(x, &n)
		//fmt.Println(diff)
		if cmp.Equal(x, &n) {
			return true
		}
	}
	return false
}

// returns true if x is found in a
func containsLink(a []dto.Link, x dto.Link) bool {
	for _, n := range a {
		if reflect.DeepEqual(x, n) {
			return true
		}
	}
	return false
}

// // returns true if x is found in a
// func contains(a []interface{}, x interface{}) bool {
// 	for _, n := range a {
// 		if reflect.DeepEqual(x, n) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func updatePatientData() (map[string]*PatientData, error) {
// 	patientMap := make(map[string]*PatientData)

// 	patientMap[TP_ROBERT_DEANGELIS] = &PatientData{nil, getFilenameFromName(TP_ROBERT_DEANGELIS), nil}

// 	for _, p := range patientMap {
// 		jsonFile, err := os.Open(patientDir + p.File)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not open test case file %s, %s", p.File, err)
// 		}
// 		defer jsonFile.Close()

// 		byteValue, err := io.ReadAll(jsonFile)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not read test case file %s, %s", p.File, err)
// 		}

// 		var patient models.Patient
// 		json.Unmarshal(byteValue, &patient)

// 		p.Data = &patient

// 		p.Id, err = svc.Create(*p.Data)
// 		if err != nil {
// 			return patientMap, fmt.Errorf("error creating patient, %s", err)
// 		}
// 	}

// 	return patientMap, nil
// }
