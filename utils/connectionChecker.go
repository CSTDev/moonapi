package utils

import (
	"errors"

	"gopkg.in/headzoo/surf.v1"
)

const baseURL string = "https://moonboard.com/"
const loginURL = "Account/Login"

func CheckConnection() (bool, error) {
	bow := surf.NewBrowser()
	err := bow.Open(baseURL + loginURL)
	if err != nil {
		return false, err
	}

	if bow.StatusCode() != 200 {
		return false, errors.New("unable to reach page")
	}
	return true, nil
}
