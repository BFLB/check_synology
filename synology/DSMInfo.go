// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Library to access Synology DSM API
//

package synology

import (
	"encoding/json"
	"errors"
	"net/url"
	//"net/http/httputil"
)

type dsmInfoResponse struct {
	Data    DSMInfo `json:"data"`
	Success bool    `json:"success"`
}

type DSMInfo struct {
	Codepage        string `json:"codepage"`
	Model           string `json:"model"`
	RAM             int    `json:"ram"`
	Serial          string `json:"serial"`
	Temperature     int    `json:"temperature"`
	TemperatureWarn bool   `json:"temperature_warn"`
	Time            string `json:"time"`
	Uptime          int    `json:"uptime"`
	Version         string `json:"version"`
	VersionString   string `json:"version_string"`
}

func (api *Syno) DSMInfo() (*DSMInfo, error) {

	// Set URL parameters
	parameters := url.Values{}
	parameters.Add("api", "SYNO.DSM.Info")
	parameters.Add("version", "2")
	parameters.Add("method", "getinfo")

	response, err := api.Get("entry.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	var payload *dsmInfoResponse
	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == false {
		return nil, errors.New(string(response))
	}

	return &payload.Data, nil
}
