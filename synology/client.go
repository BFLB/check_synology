// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Library to access Synology DSM API
//

package synology

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	//"fmt"
	//"net/http/httputil"
)

// TODO Implement error code display
/*
100 Unknown error
101 Invalid parameter
102 The requested API does not exist
103 The requested method does not exist
104 The requested version does not support the functionality
105 The logged in session does not have permission
106 Session timeout
107 Session interrupted by duplicate login
*/

var (
	ErrLoginFirst  = errors.New("login first")
	minimalVersion = 5
)

type Syno struct {
	client    *http.Client
	baseURL   string
	apiURL    string
	sid       string
	xsrfToken string
}

type DSMResponse struct {
	Data interface {
	} `json:"data"`
	Success bool `json:"success"`
}

// Initializes the http client
func (api *Syno) init(host string, port string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cj, _ := cookiejar.New(nil)

	api.client = &http.Client{
		Transport: tr,
		Jar:       cj,
	}

	if port == "443" { // FIXME
		api.baseURL = "https://" + host + "/"
	} else {
		api.baseURL = "https://" + host + ":" + port + "/"
	}
	api.apiURL = api.baseURL + ""
}

type login struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (l *login) init(user string, pass string) {
	l.Username = user
	l.Password = pass
}

// Initializes a session.
func Login(user string, pass string, host string, port string) (*Syno, error) {

	api := new(Syno)
	api.init(host, port)

	var Url *url.URL
	Url, err := url.Parse(api.baseURL)
	if err != nil {
		return nil, err
	}

	Url.Path += "webapi/auth.cgi"
	parameters := url.Values{}
	parameters.Add("api", "SYNO.API.Auth")
	parameters.Add("version", "3") // FIXME
	parameters.Add("method", "login")
	parameters.Add("account", user)
	parameters.Add("passwd", pass)
	parameters.Add("session", "FileStation")
	parameters.Add("format", "cookie")
	Url.RawQuery = parameters.Encode()

	resp, err := api.client.Get(Url.String())
	if err != nil {
		return api, err
	}

	// FIXME
	//Store the XSRF-TOKEN to be used in further requests. TODO: Get token value from cookie jar
	for _, c := range resp.Cookies() {
		if c.Name == "X-SYNO-TOKEN" {
			api.xsrfToken = c.Value
		}
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//bodyString := string(bodyBytes)

	type response struct {
		Data struct {
			Sid string `json:"sid"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	var r response
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		return nil, err
	}

	if r.Success == false {
		return nil, errors.New(string(bodyBytes))
	}

	// Login cucessful
	api.sid = r.Data.Sid
	return api, nil
}

// Initializes a session.
func (api *Syno) Get(path string, params *url.Values) ([]byte, error) {

	var Url *url.URL
	Url, err := url.Parse(api.baseURL)
	if err != nil {
		return nil, err
	}

	Url.Path += "webapi/" + path
	Url.RawQuery = params.Encode()

	resp, err := api.client.Get(Url.String())
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	//bodyString := string(bodyBytes)

	/*
		For Debugging

		fmt.Println("Dump Request")
		output, err := httputil.DumpRequest(resp.Request, true)
		if err != nil {
			fmt.Println("Error dumping request:", err)
			return nil, err
		}
		fmt.Println(string(output))

		fmt.Println("Dump Response")
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%q", dump)
	*/

	return body, nil
}

// Terminates a session
func (api *Syno) Logout(user string, pass string) error {

	var Url *url.URL
	Url, err := url.Parse(api.baseURL)
	if err != nil {
		return err
	}

	Url.Path += "webapi/auth.cgi"
	parameters := url.Values{}
	parameters.Add("api", "SYNO.API.Auth")
	parameters.Add("version", "3") // FIXME
	parameters.Add("method", "logout")
	parameters.Add("account", user)
	parameters.Add("passwd", pass)
	parameters.Add("session", "FileStation")
	parameters.Add("format", "cookie")
	Url.RawQuery = parameters.Encode()

	_, err = api.client.Get(Url.String())

	//fmt.Printf("Logout:%s", resp)

	return err

}

func (api *Syno) get(cmd string) ([]byte, error) {

	url := api.apiURL + cmd

	resp, err := api.client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *Syno) Post(cmd string, payload interface{}) ([]byte, error) {

	form := url.Values{}
	form.Set("api", "SYNO.FileIndexing.Status")
	form.Add("method", "get")
	form.Add("version", "1")

	req, err := http.NewRequest("POST", api.baseURL+"webapi/entry.cgi", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Add("X-SYNO-TOKEN", api.xsrfToken)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := api.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	// Debug output
	/*
		fmt.Println("Dump Request")
		output, err := httputil.DumpRequest(resp.Request, true)
		if err != nil {
			fmt.Println("Error dumping request:", err)
			return nil, err
		}
		fmt.Println(string(output))

		fmt.Println("Dump Response")
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%q", dump)
	*/

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *Syno) PostRaw(cmd string, payload []byte) ([]byte, error) {

	url := api.apiURL + cmd

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}

	// Set Headers
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("X-3CX-Version", "15.5.13103.5") // TODO Read from Login response
	// Add X-XSRF-TOKEN (Cross Site Forgery protection) TODO: Get it from cookie jar
	//req.Header.Add("X-XSRF-TOKEN", api.xsrfToken)

	resp, err := api.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (api *Syno) Parse(cmd string, payload interface{}, v interface{}) error {
	var body []byte
	var err error
	if payload == nil {
		body, err = api.get(cmd)
	} else {
		body, err = api.Post(cmd, payload)
	}

	if err != nil {
		return err
	}

	// Sometimes json content is surrounded with brackets. Remove them for parsing
	body = bytes.TrimLeft(body, "[")
	body = bytes.TrimRight(body, "]")

	// Quick-fix to add an empty json document in case of absence. Weired, but happens in some conditions
	if len(body) == 0 {
		s := "{}"
		body = append(body, s...)
	}

	if err := json.Unmarshal(body, &v); err != nil {

		return err
	}

	return nil
}

func (api *Syno) Save(id int) error {
	idString := strconv.Itoa(id)
	var body json.RawMessage
	body, err := api.PostRaw("edit/save", []byte(idString))

	if err != nil {
		return err
	}
	if string(body) != "{}" {
		return errors.New(string(body))
	}

	return nil
}

func (api *Syno) Cancel(id int) error {
	idString := strconv.Itoa(id)
	_, err := api.PostRaw("edit/cancel", []byte(idString))

	return err
}
