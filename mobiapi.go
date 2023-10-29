package mobiapi

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MobiAPI struct {
	client   http.Client
	domain   string
	name     string
	uid      int
	signedin bool
}

func (api *MobiAPI) request(method string, path string, body string) (*http.Response, *goquery.Document, error) {
	var resp *http.Response
	var doc *goquery.Document
	purl, err := url.Parse("https://" + api.domain + "/dziennik/" + path)
	if err != nil {
		return resp, doc, err
	}
	req, err := http.NewRequest(method, purl.String(), strings.NewReader(body))
	if err != nil {
		return resp, doc, err
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err = api.client.Do(req)
	if err != nil {
		return resp, doc, err
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return resp, doc, err
	}
	_, api.signedin = doc.Find("body").Attr("uid")

	return resp, doc, nil
}

// Check if domain is accessible and if it is, use it. If the domain doesn't contain any dots it will be treated as a mobidziennik.pl subdomain.
func (api *MobiAPI) SetDomain(domain string) error {
	if !strings.Contains(domain, ".") {
		domain = domain + ".mobidziennik.pl"
	}
	api.domain = domain
	resp, _, err := api.request("GET", "", "")
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Inaccessible")
	}
	return nil
}

// Get user's name
func (api *MobiAPI) GetName() string {
	return api.name
}

// Get UID
func (api *MobiAPI) GetUID() int {
	return api.uid
}

// Set proxy server to use and whenever to allow invalid TLS certificates. Useful for development
func (api *MobiAPI) SetupProxy(proxyurl string, noverifytls bool) error {
	if proxyurl == "" {
		api.client.Transport = &http.Transport{}
	} else {
		parsedurl, err := url.Parse(proxyurl)
		if err != nil {
			return err
		}
		api.client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(parsedurl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: noverifytls},
		}
		if _, err := api.client.Get("https://mobidziennik.pl"); err != nil {
			api.client.Transport = &http.Transport{}
			return err
		}
	}
	return nil
}

// Logout from mobiDziennik
func (api *MobiAPI) Logout() (bool, error) {
	resp, _, err := api.request("GET", "wyloguj", "")
	if err != nil {
		return false, err
	}
	if resp.Request.Response.StatusCode == 302 {
		api = nil
		api.signedin = false
		return true, nil
	}
	return false, nil
}

// Create new instance of MobiAPI.
func New(domain string) (*MobiAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	api := &MobiAPI{
		client: http.Client{
			Jar: jar,
		},
	}
	if len(domain) > 0 {
		if err := api.SetDomain(domain); err != nil {
			return api, err
		}
	}
	return api, nil
}
