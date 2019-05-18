package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/cstdev/moonapi"
	"github.com/cstdev/moonapi/utils"
)

var filePath = "./request.token"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func login(username string, password string) moonapi.MoonBoard {
	var moonBoardSession = moonapi.MoonBoard{}

	fmt.Printf("Hello %s \n", username)
	err := moonBoardSession.Login(username, password)
	check(err)

	fmt.Printf("%+v\n", moonBoardSession)

	jsonOut, err := json.Marshal(moonBoardSession.Auth)
	check(err)

	err = ioutil.WriteFile(filePath, jsonOut, 0644)
	check(err)
	return moonBoardSession
}

func reuseSession() moonapi.MoonBoard {
	// For testing so I don't actually log in each time.
	tokens, err := ioutil.ReadFile(filePath)

	var testAuth []moonapi.AuthToken

	err = json.Unmarshal([]byte(tokens), &testAuth)
	check(err)
	var moonBoardSession = moonapi.MoonBoard{}
	moonBoardSession.SetAuth(testAuth)

	fmt.Printf("%+v\n", moonBoardSession)
	return moonBoardSession
}

func main() {
	var moonBoardSession = moonapi.MoonBoard{}

	var shouldLogin = flag.Bool("login", false, "Whether to log in or use cached credentials.")
	var username = flag.String("user", "", "Enter a username to log in with.")
	var password = flag.String("pass", "", "Enter a password to log in with.")

	var order = flag.String("o", "", "Order to sort problems by: New, Grade, Rating, Repeats.")
	var desc = flag.String("d", "true", "Sort by descending.")
	var configuration = flag.String("c", "", "Board configuration: Forty, Twenty")
	var holdSet = flag.String("hs", "", "Hold Set types to include split by comma: OS, Wood, A, B, C. (default all)")
	var filter = flag.String("f", "", "Filter to apply to problems: Benchmarks, Setbyme, Myascents")
	var minGrade = flag.String("min", "", "Mininum grade to return.")
	var maxGrade = flag.String("max", "", "Maximum grade to return.")
	var page = flag.String("p", "", "Page number")
	var pageSize = flag.String("ps", "", "Page size")

	flag.Parse()

	if *shouldLogin {
		moonBoardSession = login(*username, *password)
	} else {
		moonBoardSession = reuseSession()
	}

	reqQuery := &utils.RequestQuery{
		Order:         *order,
		Asc:           *desc,
		Configuration: *configuration,
		HoldSet:       *holdSet,
		Filter:        *filter,
		MinGrade:      *minGrade,
		MaxGrade:      *maxGrade,
		Page:          *page,
		PageSize:      *pageSize,
	}

	query, err := reqQuery.Query()
	check(err)

	fmt.Printf("%+v\n", query)

	problems, err := moonBoardSession.GetProblems(query)
	check(err)

	fmt.Printf("\n\n Number of Problems: %d\n\n", problems.Total)
	fmt.Println(moonapi.ProblemsAsJSON(problems.Data))

}
