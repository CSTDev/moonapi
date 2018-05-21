package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	. "github.com/cstdev/moonapi"
	"github.com/cstdev/moonapi/query"
)

var filePath = "./request.token"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func login(username string, password string) MoonBoard {
	var moonBoardSession = MoonBoard{}

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

func reuseSession() MoonBoard {
	// For testing so I don't actually log in each time.
	tokens, err := ioutil.ReadFile(filePath)

	var testAuth []AuthToken

	err = json.Unmarshal([]byte(tokens), &testAuth)
	check(err)
	var moonBoardSession = MoonBoard{
		Auth: testAuth,
	}
	fmt.Printf("%+v\n", moonBoardSession)
	return moonBoardSession
}

func main() {
	var moonBoardSession = MoonBoard{}

	var shouldLogin = flag.Bool("login", false, "Whether to log in or use cached credentials.")
	var username = flag.String("user", "", "Enter a username to log in with.")
	var password = flag.String("pass", "", "Enter a password to log in with.")

	var order = flag.String("o", "", "Order to sort problems by: New, Grade, Rating, Repeats.")
	var desc = flag.Bool("d", true, "Sort by descending.")
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

	builder := query.New()

	if *order != "" {
		orderType, err := query.ToOrder(*order)
		check(err)
		builder.Sort(*orderType, *desc)
	}

	if *configuration != "" {
		configType, err := query.ToConfiguration(*configuration)
		check(err)

		builder.Configuration(*configType)
	}

	if *holdSet != "" {
		sets := strings.Split(*holdSet, ",")
		for _, set := range sets {
			holdType, err := query.ToHoldSet(set)
			check(err)
			builder.HoldSet(*holdType)
			fmt.Printf("Added hold set: %s", set)
		}
	}

	if *filter != "" {
		filterType, err := query.ToFilter(*filter)
		check(err)
		builder.Filter(*filterType)
	}

	if *minGrade != "" {
		minGradeType, err := query.ToGrade(*minGrade)
		check(err)

		builder.MinGrade(*minGradeType)
	}

	if *maxGrade != "" {
		maxGradeType, err := query.ToGrade(*maxGrade)
		check(err)

		builder.MaxGrade(*maxGradeType)
	}

	if *page != "" {
		intPage, err := strconv.Atoi(*page)
		check(err)
		builder.Page(intPage)
	}

	if *pageSize != "" {
		intPageSize, err := strconv.Atoi(*pageSize)
		check(err)
		builder.PageSize(intPageSize)
	}

	query, _ := builder.Build()
	fmt.Printf("%+v\n", query)

	problems, err := moonBoardSession.GetProblems(query)
	check(err)

	fmt.Printf("\n\n Number of Problems: %d\n\n", problems.Total)
	fmt.Println(ProblemsAsJSON(problems.Data))

}
