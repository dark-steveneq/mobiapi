package mobiapi

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func (api *MobiAPI) postlogin(doc *goquery.Document) {
	if val, exists := doc.Find("body").Attr("uid"); exists {
		api.uid, _ = strconv.Atoi(val)
		api.name = doc.Find("#botton div strong").Text()
	}
}

// Authenticate with password.
func (api *MobiAPI) PasswordAuth(login, password string) (bool, error) {
	if login == "" || password == "" {
		return false, errors.New("EmptyCred")
	}
	purl, err := url.Parse("https://" + api.domain)
	if err != nil {
		return false, err
	}
	api.client.Jar.SetCookies(purl, []*http.Cookie{})

	resp, doc, err := api.request("POST", "logowanie", "login="+login+"&haslo="+password)
	if err != nil {
		return false, err
	}

	if resp.Request.Response.StatusCode == 302 {
		api.signedin = true
		api.postlogin(doc)
		return true, nil
	}

	return false, nil
}

// Authenticate with provided token from cookie
func (api *MobiAPI) TokenAuth(token string) (bool, error) {
	purl, err := url.Parse("https://" + api.domain)
	if err != nil {
		return false, err
	}
	api.client.Jar.SetCookies(purl, []*http.Cookie{})

	for _, cookie := range api.client.Jar.Cookies(purl) {
		if cookie.Name != "SERVERID" {
			cookie.Value = token
			api.client.Jar.SetCookies(purl, []*http.Cookie{cookie})
			break
		}
	}

	resp, doc, err := api.request("GET", "logowanie", "")
	if err != nil {
		return false, err
	}

	if resp.Request.URL.Path == "/dziennik/historia" {
		api.signedin = true
		api.postlogin(doc)
		return true, nil
	} else {
		return false, errors.New("AuthUnable")
	}
}
