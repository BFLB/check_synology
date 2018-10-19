// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Library to access Synology DSM API

package synology

import (
	"encoding/json"
	"errors"
	//"fmt"
	"net/url"
)

type dualEnclosureResponse struct {
	Data    DualEnclosure `json:"data"`
	Error struct {
		Code int `json:"code"`
	} `json:"error, omitempty"`
	Success bool          `json:"success"`
}

type DualEnclosure struct {
	AHAInfo struct {
		EncCnt     int `json:"enc_cnt"`
		Enclosures []struct {
			Disks       []int  `json:"disks"`
			Fans        []int  `json:"fans"`
			Links       []int  `json:"links"`
			MaxDisk     int    `json:"max_disk"`
			Model       string `json:"model"`
			ModelID     int    `json:"model_id"`
			Powers      []Power  `json:"powers"` // Values: 1=normal 0=anormal
			SN          string `json:"sn"`
			Temperature int    `json:"temperature"`
		} `json:"enclosures"`
		HostCnt int `json:"host_cnt"`
		Hosts   []struct {
			Fans        []int  `json:"fans"`
			HostType    int    `json:"host_type"`
			Links       []int  `json:"links"`
			Modelname   string `json:"modelname"`
			Powers      []Power  `json:"powers"`
			SN          string `json:"sn"`
			Temperature int    `json:"temperature"`
		} `json:"hosts"`
	} `json:"AHAInfo"`
}

type Power int

func (enc Power) String() string {
	switch enc {
	case 1:
		return "normal"
	default:
		return "anormal"
	}
}

func (api *Syno) DualEnclosure() (*DualEnclosure, error) {
	// Set URL parameters
	parameters := url.Values{}
	parameters.Add("api", "SYNO.Storage.CGI.DualEnclosure")
	parameters.Add("version", "1")
	parameters.Add("method", "load")

	response, err := api.Get("entry.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	var payload *dualEnclosureResponse
	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == false {
		return nil, errors.New(string(response))
	}

	return &payload.Data, nil

}
