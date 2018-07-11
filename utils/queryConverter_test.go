package utils

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error returned: %v", err)
		t.FailNow()
	}
}

func compare(actual string, expected string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected:\n %s \n Recieved: \n %s", expected, actual)
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	retCode := m.Run()
	os.Exit(retCode)
}

func TestOrderIsAddedToTheQuery(t *testing.T) {
	req := &RequestQuery{
		Order: "New",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "New-desc"
	compare(query.Sort(), expected, t)
}

func TestOrderAscIsAddedToTheQuery(t *testing.T) {
	req := &RequestQuery{
		Order: "Grade",
		Asc:   "true",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "GradeAsc-asc"
	compare(query.Sort(), expected, t)
}

func TestConfigurationIsAddedToTheQuery(t *testing.T) {
	req := &RequestQuery{
		Configuration: "Forty",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "Configuration~eq~'40° MoonBoard'~and~MinGrade~eq~'6A+'~and~MaxGrade~eq~'8B+'"
	compare(query.Filter(), expected, t)
}

func TestHoldSetIsAddedToTheQuery(t *testing.T) {
	req := &RequestQuery{
		HoldSet: "A",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "Holdsets~eq~'hold set a'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	compare(query.Filter(), expected, t)
}

func TestMultipleHoldSetsAddedToQuery(t *testing.T) {

	req := &RequestQuery{
		HoldSet: "A, os",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "Holdsets~eq~'hold set a,original school holds'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	compare(query.Filter(), expected, t)
}

func TestFilterAddedToQuery(t *testing.T) {
	req := &RequestQuery{
		Filter: "Benchmarks",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "Benchmarks~eq~''~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	compare(query.Filter(), expected, t)
}

func TestMinGradeAddedToQuery(t *testing.T) {
	req := &RequestQuery{
		MinGrade: "7A+",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "MinGrade~eq~'7A+'~and~MaxGrade~eq~'8B+'"
	compare(query.Filter(), expected, t)
}

func TestMaxGradeAddedToQuery(t *testing.T) {
	req := &RequestQuery{
		MaxGrade: "7a+",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := "MinGrade~eq~'5+'~and~MaxGrade~eq~'7A+'"
	compare(query.Filter(), expected, t)
}

func TestPageIsAddedToQuery(t *testing.T) {
	req := &RequestQuery{
		Page: "4",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := 4
	if query.Page() != expected {
		t.Errorf("Expected:\n %d \n Recieved: \n %d", expected, query.Page())
		t.FailNow()
	}
}

func TestPageSizeIsAddedToQuery(t *testing.T) {
	req := &RequestQuery{
		Page: "50",
	}

	query, err := req.Query()
	checkError(t, err)

	expected := 50
	if query.Page() != expected {
		t.Errorf("Expected:\n %d \n Recieved: \n %d", expected, query.Page())
		t.FailNow()
	}
}
