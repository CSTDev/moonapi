package query

import (
	"testing"
)

func TestBuilderSortNewest(t *testing.T) {
	expected := "New-desc"
	builder := New()
	query, _ := builder.Sort(Newest, false).Build()
	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}
}

func TestBuilderSortDifficulty(t *testing.T) {
	expected := "GradeAsc-asc"
	builder := New()
	query, err := builder.Sort(Difficulty, true).Build()

	if err != nil {
		t.Errorf("Unexpected error. %s", err[0].Error())
		t.FailNow()
	}

	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}

	builder = New()
	expected = "GradeDesc-desc"
	query, err = builder.Sort(Difficulty, false).Build()

	if err != nil {
		t.Errorf("Unexpected error. %s", err[0].Error())
		t.FailNow()
	}

	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}
}

func TestBuilderSortRating(t *testing.T) {
	expected := "Rating-desc"
	builder := New()
	query, _ := builder.Sort(Rating, false).Build()
	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}
}

func TestBuilderSortRepeats(t *testing.T) {
	expected := "RepeatsAsc-asc"
	builder := New()
	query, err := builder.Sort(Repeats, true).Build()

	if err != nil {
		t.Errorf("Unexpected error. %s", err[0].Error())
		t.FailNow()
	}

	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}

	builder = New()
	expected = "RepeatsDesc-desc"
	query, err = builder.Sort(Repeats, false).Build()

	if err != nil {
		t.Errorf("Unexpected error. %s", err[0].Error())
		t.FailNow()
	}

	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
	}
}

func TestBuilderOnlyOneSort(t *testing.T) {

	expected := "New-desc"

	builder := New()
	query, err := builder.Sort(Repeats, true).Sort(Newest, false).Build()

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}

	expectedError := "Can only sort by one parameter, defaulting to the last provided."
	if err[0].Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err[0].Error(), expectedError)
	}

	if query.Sort() != expected {
		t.Errorf("Query sort was incorrect, got %s, want: %s", query.Sort(), expected)
		t.FailNow()
	}

}

func TestBuilderConfigurationForty(t *testing.T) {
	expected := "Configuration~eq~'40° MoonBoard'~and~MinGrade~eq~'6A+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Configuration(Forty).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderConfigurationTwenty(t *testing.T) {
	expected := "Configuration~eq~'25° MoonBoard'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Configuration(Twenty).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderConfigurationBoth(t *testing.T) {
	expected := "Configuration~eq~'40° MoonBoard,25° MoonBoard'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Configuration(Forty).Configuration(Twenty).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetOS(t *testing.T) {
	expected := "Holdsets~eq~'original school holds'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(OS).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetWood(t *testing.T) {
	expected := "Holdsets~eq~'wooden holds'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(Wood).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetA(t *testing.T) {
	expected := "Holdsets~eq~'hold set a'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(A).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetB(t *testing.T) {
	expected := "Holdsets~eq~'hold set b'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(B).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetC(t *testing.T) {
	expected := "Holdsets~eq~'hold set c'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(C).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetMultiple(t *testing.T) {
	expected := "Holdsets~eq~'hold set c,wooden holds'~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(C).HoldSet(Wood).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderHoldSetAndConfiguration(t *testing.T) {
	expected := "Configuration~eq~'40° MoonBoard'~and~Holdsets~eq~'hold set c'~and~MinGrade~eq~'6A+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.HoldSet(C).Configuration(Forty).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderFilterMyAscents(t *testing.T) {
	expected := "Myascents~eq~''~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Filter(MyAscents).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderFilterBenchmarks(t *testing.T) {
	expected := "Benchmarks~eq~''~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Filter(Benchmarks).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderFilterSetByMe(t *testing.T) {
	expected := "Setbyme~eq~''~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Filter(SetByMe).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderGradeMin(t *testing.T) {
	expected := "MinGrade~eq~'7A+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.MinGrade(SevenAPlus).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderGradeMax(t *testing.T) {
	expected := "MinGrade~eq~'5+'~and~MaxGrade~eq~'7B'"
	builder := New()
	query, _ := builder.MaxGrade(SevenB).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestBuilderMaxGradeLowerThanMinErrors(t *testing.T) {
	expectedError := "Min grade cannot be higher than max grade."
	builder := New()

	_, err := builder.MinGrade(EightA).MaxGrade(SixC).Build()

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}

	if err[0].Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err[0].Error(), expectedError)
	}
}

func TestBuilderFilterMultiple(t *testing.T) {
	expected := "Myascents~eq~''~and~Benchmarks~eq~''~and~MinGrade~eq~'5+'~and~MaxGrade~eq~'8B+'"
	builder := New()
	query, _ := builder.Filter(MyAscents).Filter(Benchmarks).Build()
	if query.Filter() != expected {
		t.Errorf("Query Filter was incorrect, got %s, want: %s", query.Filter(), expected)
	}
}

func TestAddPageNumber(t *testing.T) {
	expected := 2
	builder := New()

	query, err := builder.MinGrade(SevenA).Page(2).Build()
	if err != nil {
		t.Errorf("Unexpected error recieved")
		t.FailNow()
	}

	if query.Page() != expected {
		t.Errorf("Page was incorrect. Got %d\n Expected %d", query.Page(), expected)
	}
}

func TestBadPageNumber(t *testing.T) {
	expectedError := "Page number cannot be below 1."
	builder := New()

	query, err := builder.MinGrade(SevenA).Page(0).Build()
	if err == nil {
		t.Errorf("Eexpected error not recieved")
		t.FailNow()
	}

	if err[0].Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s\n Expected: %s", err[0].Error(), expectedError)
	}

	expected := 1
	if query.Page() != expected {
		t.Errorf("Page was incorrect. Got %d\n Expected %d", query.Page(), expected)
	}
}

func TestSetPageSize(t *testing.T) {
	expected := 50
	builder := New()

	query, err := builder.MaxGrade(SevenB).PageSize(50).Build()
	if err != nil {
		t.Errorf("Unexpected error recieved")
		t.FailNow()
	}

	if query.PageSize() != expected {
		t.Errorf("Page Size was incorrect. Got %d\n Expected %d", query.PageSize(), expected)
	}
}

func TestSetPageSizeTooLargeOrSmall(t *testing.T) {
	expectedError := "Page size must be between 1 and 100"
	builder := New()

	query, err := builder.MaxGrade(SevenB).PageSize(150).Build()
	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}

	if err[0].Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s\n Expected: %s", err[0].Error(), expectedError)
	}

	expected := 15
	if query.PageSize() != expected {
		t.Errorf("Page Size was incorrect. Got %d\n Expected %d", query.PageSize(), expected)
	}

	query, err = builder.MaxGrade(SevenB).PageSize(0).Build()
	if len(err) < 2 {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}

	if err[1].Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s\n Expected: %s", err[1].Error(), expectedError)
	}

	if query.PageSize() != expected {
		t.Errorf("Page Size was incorrect. Got %d\n Expected %d", query.PageSize(), expected)
	}
}

func TestToOrderErrorsOnInvalidValue(t *testing.T) {
	_, err := ToOrder("Test")
	if err == nil {
		t.Error("Expected error not recieved")
	}
}

func TestToConfigurationErrorsOnInvalidValue(t *testing.T) {
	_, err := ToConfiguration("Test")
	if err == nil {
		t.Error("Expected error not recieved")
	}
}

func TestToHoldSetErrorsOnInvalidValue(t *testing.T) {
	_, err := ToHoldSet("Test")
	if err == nil {
		t.Error("Expected error not recieved")
	}
}

func TestToFilterErrorsOnInvalidValue(t *testing.T) {
	_, err := ToHoldSet("Test")
	if err == nil {
		t.Error("Expected error not recieved")
	}
}

func TestToGradeErrorsOnInvalidValue(t *testing.T) {
	_, err := ToGrade("Test")
	if err == nil {
		t.Error("Expected error not recieved")
	}
}
