package moonapi

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cstdev/moonapi/query"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestSuccessfulLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://moonboard.com/Account/Login",
		httpmock.NewStringResponder(200, loginForm))

	httpmock.RegisterResponder("POST", "https://moonboard.com/Account/Login",
		func(req *http.Request) (*http.Response, error) {

			respRecorder := httptest.NewRecorder()
			handler := func(w http.ResponseWriter, r *http.Request) {

				cookie := http.Cookie{Name: "__RequestVerificationToken", Value: "Value1"}
				cookie2 := http.Cookie{Name: "_MoonBoard", Value: "Value2"}

				http.SetCookie(w, &cookie)
				http.SetCookie(w, &cookie2)

				io.WriteString(w, loginForm)
			}

			handler(respRecorder, req)
			resp := respRecorder.Result()
			resp.Request = req // Important
			return resp, nil

		},
	)

	var session = MoonBoard{}
	err := session.Login("TestUser", "Password1")
	if err != nil {
		t.Errorf("Expected to login, recieved error: %s", err.Error())
	}
}

func TestUnableToLoadLoginPage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://moonboard.com/Account/Login",
		httpmock.NewStringResponder(500, badLoginForm))

	var session = MoonBoard{}
	err := session.Login("TestUser", "Password1")

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}

}

func TestUnableToSubmitLoginForm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://moonboard.com/Account/Login",
		httpmock.NewStringResponder(200, loginForm))

	httpmock.RegisterResponder("POST", "https://moonboard.com/Account/Login",
		httpmock.NewErrorResponder(errors.New("Failed to submit login form")))

	var session = MoonBoard{}
	err := session.Login("TestUser", "Password1")

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "Failed to submit log-in"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}
}

func TestInvalidLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://moonboard.com/Account/Login",
		httpmock.NewStringResponder(200, loginForm))

	httpmock.RegisterResponder("POST", "https://moonboard.com/Account/Login",
		httpmock.NewStringResponder(200, loginForm))

	var session = MoonBoard{}
	err := session.Login("TestUser", "Password1")

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "failed to log-in, moonboard cookie not returned"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}

}

const loginForm string = `<body><Form action="/Account/Login" method="post" id="frmLogin"></Form></body>`

const badLoginForm string = `<body><h1>Unable to load login page</h1></body>`

func TestErrorOnSessionMissingAuthTokens(t *testing.T) {
	var session = MoonBoard{}
	builder := query.New()
	q, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()

	_, err := session.GetProblems(q)

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "Required _Moonboard or __RequestVerificationToken Auth Tokens missing"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}
}

func TestErrorOnSessionWrongAuthTokens(t *testing.T) {

	var testAuth []AuthToken
	testAuth = append(testAuth, AuthToken{Name: "__Token", Value: "RequestToken"})
	testAuth = append(testAuth, *testReqCookie)
	var session = MoonBoard{
		Auth: testAuth,
	}

	builder := query.New()
	q, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()

	_, err := session.GetProblems(q)

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "Required _Moonboard or __RequestVerificationToken Auth Tokens missing"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}
}

func TestErrorOnSessionExpired(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://moonboard.com/Problems/GetProblems",
		func(req *http.Request) (*http.Response, error) {

			respRecorder := httptest.NewRecorder()
			handler := func(w http.ResponseWriter, r *http.Request) {

				cookie := http.Cookie{Name: "__RequestVerificationToken", Value: "Value1"}
				cookie2 := http.Cookie{Name: "_MoonBoard", Value: "Value2"}

				http.SetCookie(w, &cookie)
				http.SetCookie(w, &cookie2)

				io.WriteString(w, loginForm)
			}

			handler(respRecorder, req)
			resp := respRecorder.Result()
			req.URL, _ = url.Parse("https://moonboard.com/Account/Login")
			resp.Request = req // Important
			return resp, nil

		},
	)

	var testAuth []AuthToken
	testAuth = append(testAuth, *testMoonCookie)
	testAuth = append(testAuth, *testReqCookie)
	var session = MoonBoard{
		Auth: testAuth,
	}

	builder := query.New()
	q, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()

	_, err := session.GetProblems(q)

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "session expired, please log in"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}
}

func TestValidGetProbelmsQueryReturnsProblems(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://moonboard.com/Problems/GetProblems",
		httpmock.NewStringResponder(200, problems))

	var testAuth []AuthToken
	testAuth = append(testAuth, *testMoonCookie)
	testAuth = append(testAuth, *testReqCookie)
	var session = MoonBoard{
		Auth: testAuth,
	}

	builder := query.New()
	q, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()

	problems, err := session.GetProblems(q)

	if err != nil {
		t.Errorf("Error recieved: %v", err)
		t.FailNow()
	}

	expected := 2
	if len(problems.Data) != expected || problems.Total != expected {
		t.Errorf("Expected there to be %d problems, \n there were %d\n The total was %d", expected, len(problems.Data), problems.Total)
	}
}

func TestQueryParametersArePassedToEnpoint(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expectedString := `filter=Configuration~eq~%2740%C2%B0+MoonBoard%27~and~Benchmarks~eq~%27%27~and~MinGrade~eq~%276A%2B%27~and~MaxGrade~eq~%277B%27&group=&page=4&pageSize=15&sort=`

	httpmock.RegisterResponder("POST", "https://moonboard.com/Problems/GetProblems",
		func(req *http.Request) (*http.Response, error) {
			body := req.Body
			buf := new(bytes.Buffer)
			buf.ReadFrom(body)
			s := buf.String()

			if s != expectedString {
				t.Errorf("Request did not contain expected body. Got %s\n Expected %s", s, expectedString)
				t.FailNow()
			}

			respRecorder := httptest.NewRecorder()
			handler := func(w http.ResponseWriter, r *http.Request) {

				cookie := http.Cookie{Name: "__RequestVerificationToken", Value: "Value1"}
				cookie2 := http.Cookie{Name: "_MoonBoard", Value: "Value2"}

				http.SetCookie(w, &cookie)
				http.SetCookie(w, &cookie2)

				io.WriteString(w, problems)
			}

			handler(respRecorder, req)
			resp := respRecorder.Result()
			resp.Request = req // Important
			return resp, nil

		},
	)

	var testAuth []AuthToken
	testAuth = append(testAuth, *testMoonCookie)
	testAuth = append(testAuth, *testReqCookie)
	var session = MoonBoard{
		Auth: testAuth,
	}

	builder := query.New()
	q, _ := builder.Configuration(query.Forty).Filter(query.Benchmarks).Page(4).PageSize(15).MinGrade(query.SixAPlus).MaxGrade(query.SevenB).Build()

	session.GetProblems(q)

}

func TestNonTwoHundredResponseReturnsError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://moonboard.com/Problems/GetProblems",
		httpmock.NewStringResponder(500, "<h1>Internal server error</h1>"))

	var testAuth []AuthToken
	testAuth = append(testAuth, *testMoonCookie)
	testAuth = append(testAuth, *testReqCookie)
	var session = MoonBoard{
		Auth: testAuth,
	}

	builder := query.New()
	q, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()

	_, err := session.GetProblems(q)

	if err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
	expectedError := "Server returned error status: 500"
	if err.Error() != expectedError {
		t.Errorf("Incorrect error provided. Got: %s Expected: %s", err.Error(), expectedError)
	}
}

var testReqCookie = &AuthToken{
	Name:  "__RequestVerificationToken",
	Value: "RequestToken",
}

var testMoonCookie = &AuthToken{
	Name:  "_MoonBoard",
	Value: "MoonToken",
}

var problems = `{"Data":[{"Method":"Feet follow hands","Name":"SOFT WOOD RH","Grade":"7B","UserGrade":"7B","MoonBoardConfiguration":{"Id":1,"Description":"40° MoonBoard","LowGrade":null,"HighGrade":null},"MoonBoardConfigurationId":0,"Setter":{"Id":"5FC09F63-05F3-4DAE-A1A5-3AC22C37139A","Nickname":"Ben Moon","Firstname":"Ben","Lastname":"Moon","City":"Sheffield","Country":"United Kingdom","ProfileImageUrl":"/Content/Account/Images/default-profile.png?636608804534169534","CanShareData":true},"FirstAscender":false,"Rating":0,"UserRating":3,"Repeats":19,"Attempts":0,"Holdsetup":{"Id":15,"Description":"MoonBoard Masters 2017","Setby":null,"DateInserted":null,"DateUpdated":null,"DateDeleted":null,"IsLocked":false,"Holdsets":null,"MoonBoardConfigurations":null,"HoldLayoutId":0,"AllowClimbMethods":true},"IsBenchmark":true,"Moves":[{"Id":1729512,"Description":"C18","IsStart":false,"IsEnd":true},{"Id":1729513,"Description":"D16","IsStart":false,"IsEnd":false},{"Id":1729514,"Description":"G5","IsStart":true,"IsEnd":false},{"Id":1729515,"Description":"G12","IsStart":false,"IsEnd":false},{"Id":1729516,"Description":"G15","IsStart":false,"IsEnd":false},{"Id":1729517,"Description":"H7","IsStart":false,"IsEnd":false},{"Id":1729518,"Description":"J5","IsStart":true,"IsEnd":false},{"Id":1729519,"Description":"K10","IsStart":false,"IsEnd":false}],"Holdsets":null,"Locations":[{"Id":0,"Holdset":null,"Description":null,"X":195,"Y":88,"Color":"0xFF0000","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":245,"Y":186,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":395,"Y":738,"Color":"0x00FF00","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":395,"Y":388,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":395,"Y":238,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":445,"Y":636,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":545,"Y":736,"Color":"0x00FF00","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":595,"Y":488,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"}],"RepeatText":"19 climbers  have repeated this problem","NumberOfTries":null,"NameForUrl":"soft-wood-rh","Id":318731,"ApiId":0,"DateInserted":"/Date(1524237072990)/","DateUpdated":null,"DateDeleted":null,"DateTimeString":"20 Apr 2018 16:11"},{"Method":"Feet follow hands","Name":"SOFT WOOD LH","Grade":"7B","UserGrade":"7B","MoonBoardConfiguration":{"Id":1,"Description":"40° MoonBoard","LowGrade":null,"HighGrade":null},"MoonBoardConfigurationId":0,"Setter":{"Id":"5FC09F63-05F3-4DAE-A1A5-3AC22C37139A","Nickname":"Ben Moon","Firstname":"Ben","Lastname":"Moon","City":"Sheffield","Country":"United Kingdom","ProfileImageUrl":"/Content/Account/Images/default-profile.png?636608804534325784","CanShareData":true},"FirstAscender":false,"Rating":0,"UserRating":3,"Repeats":19,"Attempts":0,"Holdsetup":{"Id":15,"Description":"MoonBoard Masters 2017","Setby":null,"DateInserted":null,"DateUpdated":null,"DateDeleted":null,"IsLocked":false,"Holdsets":null,"MoonBoardConfigurations":null,"HoldLayoutId":0,"AllowClimbMethods":true},"IsBenchmark":true,"Moves":[{"Id":1729520,"Description":"A10","IsStart":false,"IsEnd":false},{"Id":1729521,"Description":"B5","IsStart":true,"IsEnd":false},{"Id":1729522,"Description":"D7","IsStart":false,"IsEnd":false},{"Id":1729523,"Description":"E5","IsStart":true,"IsEnd":false},{"Id":1729524,"Description":"E12","IsStart":false,"IsEnd":false},{"Id":1729525,"Description":"E15","IsStart":false,"IsEnd":false},{"Id":1729526,"Description":"H16","IsStart":false,"IsEnd":false},{"Id":1729527,"Description":"I18","IsStart":false,"IsEnd":true}],"Holdsets":null,"Locations":[{"Id":0,"Holdset":null,"Description":null,"X":95,"Y":488,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":145,"Y":736,"Color":"0x00FF00","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":245,"Y":636,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":295,"Y":738,"Color":"0x00FF00","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":295,"Y":388,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":295,"Y":238,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":445,"Y":186,"Color":"0x0000FF","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"},{"Id":0,"Holdset":null,"Description":null,"X":495,"Y":88,"Color":"0xFF0000","Rotation":0,"Type":0,"HoldNumber":null,"Direction":0,"DirectionString":"N"}],"RepeatText":"19 climbers  have repeated this problem","NumberOfTries":null,"NameForUrl":"soft-wood-lh","Id":318730,"ApiId":0,"DateInserted":"/Date(1524237026033)/","DateUpdated":null,"DateDeleted":null,"DateTimeString":"20 Apr 2018 16:10"}], "Total":2}`

func TestProblemsToJSONReturnsJSON(t *testing.T) {
	expectedJSON := `[{"Method":"Feet follow hands","Name":"SOFT WOOD RH","Grade":"7B","UserGrade":"7B","MoonBoardConfiguration":{"Id":1,"Description":"40° MoonBoard","LowGrade":null,"HighGrade":null},"MoonBoardConfigurationId":0,"Setter":{"Id":"5FC09F63-05F3-4DAE-A1A5-3AC22C37139A","Nickname":"Ben Moon","Firstname":"Ben","Lastname":"Moon","City":"Sheffield","Country":"United Kingdom","ProfileImageUrl":"/Content/Account/Images/default-profile.png?636608804534169534","CanShareData":true},"FirstAscender":false,"Rating":0,"UserRating":3,"Repeats":19,"Attempts":0,"Holdsetup":{"Id":15,"Description":"MoonBoard Masters 2017","Setby":null,"DateInserted":null,"DateUpdated":null,"DateDeleted":null,"IsLocked":false,"Holdsets":null,"MoonBoardConfigurations":null,"HoldLayoutId":0,"AllowClimbMethods":true},"IsBenchmark":true,"Moves":null,"Holdsets":null,"Locations":null,"RepeatText":"19 climbers  have repeated this problem","NumberOfTries":null,"NameForUrl":"soft-wood-rh","Id":318731,"ApiId":0,"DateInserted":"/Date(1524237072990)/","DateUpdated":null,"DateDeleted":null,"DateTimeString":"20 Apr 2018 16:11"}]`

	var problems []Problem

	problem := &Problem{
		Method:    "Feet follow hands",
		Name:      "SOFT WOOD RH",
		Grade:     "7B",
		UserGrade: "7B",
		MoonBoardConfiguration: MoonBoardConfiguration{
			ID:          1,
			Description: "40° MoonBoard",
			LowGrade:    nil,
			HighGrade:   nil,
		},
		MoonBoardConfigurationID: 0,
		Setter: Setter{
			ID:              "5FC09F63-05F3-4DAE-A1A5-3AC22C37139A",
			Nickname:        "Ben Moon",
			Firstname:       "Ben",
			Lastname:        "Moon",
			City:            "Sheffield",
			Country:         "United Kingdom",
			ProfileImageURL: "/Content/Account/Images/default-profile.png?636608804534169534",
			CanShareData:    true,
		},
		FirstAscender: false,
		Rating:        0,
		UserRating:    3,
		Repeats:       19,
		Attempts:      0,
		Holdsetup: HoldSetup{
			ID:                      15,
			Description:             "MoonBoard Masters 2017",
			Setby:                   nil,
			DateInserted:            nil,
			DateUpdated:             nil,
			DateDeleted:             nil,
			IsLocked:                false,
			Holdsets:                nil,
			MoonBoardConfigurations: nil,
			HoldLayoutID:            0,
			AllowClimbMethods:       true,
		},
		IsBenchmark:    true,
		Holdsets:       nil,
		RepeatText:     "19 climbers  have repeated this problem",
		NumberOfTries:  nil,
		NameForURL:     "soft-wood-rh",
		ID:             318731,
		APIID:          0,
		DateInserted:   "/Date(1524237072990)/",
		DateUpdated:    nil,
		DateDeleted:    nil,
		DateTimeString: "20 Apr 2018 16:11",
	}

	problems = append(problems, *problem)

	returnedJSON, err := ProblemsAsJSON(problems)

	if err != nil {
		t.Errorf("Unexpected error recieved")
		t.FailNow()
	}

	if returnedJSON != expectedJSON {
		t.Errorf("Returned json does not match expected. \n Got %s \n Expected %s", returnedJSON, expectedJSON)
	}
}

func TestIfJSONConversionFailsErrorIsReturned(t *testing.T) {
	var problems []Problem
	problems = append(problems, Problem{})

	json, err := ProblemsAsJSON(problems)

	if json == "" && err == nil {
		t.Errorf("Expected error not recieved")
		t.FailNow()
	}
}
