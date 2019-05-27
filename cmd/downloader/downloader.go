package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cstdev/moonapi"
	"github.com/cstdev/moonapi/query"
	log "github.com/sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	logLevel := os.Getenv("LOG_LEVEL")

	log.SetFormatter(&log.JSONFormatter{})

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		break
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if username == "" {
		log.Fatal("Username environment variable set")
	}

	if password == "" {
		log.Fatal("Password environment variable set")
	}

	moonBoardSession := login(username, password)
	log.WithFields(log.Fields{
		"mbSession": moonBoardSession,
	}).Debug("Logged in session")

	builder := query.New()
	query, _ := builder.Build()
	problem, err := moonBoardSession.GetProblems(query)

	check(err)
	fmt.Printf("Problems: /n %v", problem)

}

func login(username string, password string) moonapi.MoonBoard {
	var moonBoardSession = moonapi.MoonBoard{}

	fmt.Printf("Hello %s \n", username)
	err := moonBoardSession.Login(username, password)
	check(err)

	fmt.Printf("%+v\n", moonBoardSession)

	return moonBoardSession
}
