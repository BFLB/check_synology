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

type systemStatusResponse struct {
	Data    SystemStatus `json:"data"`
	Success bool         `json:"success"`
}

type SystemStatus struct {
	IsSystemCrashed bool `json:"is_system_crashed"`
	UpgradeReady    bool `json:"upgrade_ready"`
}

func (api *Syno) SystemStatus() (*SystemStatus, error) {

	// Set URL parameters
	parameters := url.Values{}
	parameters.Add("api", "SYNO.Core.System.Status")
	parameters.Add("version", "1")
	parameters.Add("method", "get")

	response, err := api.Get("entry.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	var payload *systemStatusResponse
	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == false {
		return nil, errors.New(string(response))
	}

	return &payload.Data, nil

}
