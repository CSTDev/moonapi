package moonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	. "github.com/cstdev/moonapi/query"
	"github.com/golang/glog"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/headzoo/surf.v1"
)

// AuthToken contains the values of the cookie required to be used as authentication
// against the website.
type AuthToken struct {
	Name  string
	Value string
}

// MoonBoardApi provides methods for interacting with the website
type MoonBoardApi interface {
	Login(username string, password string) error
	GetProblems(query Query) (MbResponse, error)
}

// MoonBoard contains all the AuthTokens (cookies) required
type MoonBoard struct {
	Auth []AuthToken
}

const baseUrl string = "https://moonboard.com/"
const loginUrl = "Account/Login"
const getProblemsUrl = "Problems/GetProblems"

// Login takes a username and password, then attempts to use these to
// enter into the website's login form and submit it, storing the resulting
// cookies as AuthTokens
func (m *MoonBoard) Login(username string, password string) error {
	fmt.Printf("Hi %s\n", username)
	bow := surf.NewBrowser()
	err := bow.Open(baseUrl + loginUrl)
	if err != nil {
		glog.Info("Unable to open Login Page.")
		return err
	}

	fm, err := bow.Form("#frmLogin")
	if err != nil {
		glog.Info("Unable to find Login form.")
		return err
	}

	fm.Input("Login.Username", username)
	fm.Input("Login.Password", password)
	fm.Input("Login.RememberMe", "false")

	if fm.Submit() != nil {
		glog.Info("Failed to submit log-in")
		return errors.New("Failed to submit log-in")
	}

	var response []AuthToken
	var successResponse = false
	for _, cookie := range bow.SiteCookies() {
		token := AuthToken{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
		if cookie.Name == "_MoonBoard" {
			successResponse = true
		}
		response = append(response, token)
	}

	if !successResponse {
		//fmt.Printf("Response: %v", response)
		return errors.New("failed to log-in, moonboard cookie not returned")
	}

	m.Auth = response

	return nil

}

func tokenToCookie(token AuthToken) *http.Cookie {
	return &http.Cookie{
		Name:  token.Name,
		Value: token.Value,
	}
}

// GetProblems can be called on a session object to return all problems that
// match the provided critieria from the Query passed in.
// It requires the session to provide the
// _MoonBoard and __RequestVerificationToken AuthToken
// errors are retuned if these are missing or the session has expired.
func (m *MoonBoard) GetProblems(query Query) (MbResponse, error) {
	v := url.Values{}
	v.Set("page", strconv.Itoa(query.Page()))
	v.Add("pageSize", strconv.Itoa(query.PageSize()))
	v.Add("group", "")
	v.Add("sort", query.Sort())
	v.Add("filter", query.Filter())

	res := MbResponse{}

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	var cookies []*http.Cookie

	if len(m.Auth) == 0 {
		return res, errors.New("Required _Moonboard or __RequestVerificationToken Auth Tokens missing")
	}

	containsAuth := false
	for _, token := range m.Auth {
		if token.Name == "_MoonBoard" {
			containsAuth = true
			break
		}
	}

	if !containsAuth {
		return res, errors.New("Required _Moonboard or __RequestVerificationToken Auth Tokens missing")
	}

	cookies = append(cookies, tokenToCookie(m.Auth[0]))
	cookies = append(cookies, tokenToCookie(m.Auth[1]))
	u, _ := url.Parse(baseUrl)
	jar.SetCookies(u, cookies)
	bow := surf.NewBrowser()
	bow.SetCookieJar(jar)

	err := bow.PostForm(baseUrl+getProblemsUrl, v)

	if err != nil {
		return res, err
	}

	if strings.Contains(bow.Url().String(), "/Account/Login") {
		//fmt.Println("Session Exprired")
		return res, errors.New("session expired, please log in")
	}

	if bow.StatusCode() != 200 {
		return res, errors.New("Server returned error status: " + strconv.Itoa(bow.StatusCode()))
	}
	response := strings.Replace(bow.Body(), "&#34;", "\"", -1)
	//fmt.Printf("Response: %v \n", response)

	err = json.Unmarshal([]byte(response), &res)

	if err != nil {
		//fmt.Println("Error on unmarshal")
		return res, err
	}

	return res, nil

}

// ProblemsAsJSON takes an array of Problem and returns
// the JSON representation of those problems.
// Errors are returned if it fails to marshal the object.
func ProblemsAsJSON(problems []Problem) (string, error) {
	out, err := json.Marshal(problems)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
