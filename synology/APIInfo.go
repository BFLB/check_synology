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
	"fmt"
	"net/url"
)

type APIInfoElement struct {
	API           string
	MaxVersion    int
	MinVersion    int
	Path          string
	RequestFormat string
}

func (api *Syno) APIInfo() ([]APIInfoElement, error) {
	// Set URL parameters
	parameters := url.Values{}
	parameters.Add("api", "SYNO.API.Info")
	parameters.Add("version", "1")
	parameters.Add("method", "query")

	response, err := api.Get("query.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Response:%v\n", response) // TODO Remove

	var payload map[string]interface{}

	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Payload:%v\n", payload) // TODO Remove
	success := payload["success"].(bool)

	if success == false {
		return nil, errors.New(string(response))
	}

	// TODO Improve. Try to use unmarshal. Did not work right now.
	data := payload["data"].(map[string]interface{})
	var elements []APIInfoElement
	for k, v := range data {
		var element APIInfoElement
		var rawElement map[string]interface{}
		rawElement = v.(map[string]interface{})
		element.API = k
		element.MaxVersion = int(rawElement["maxVersion"].(float64))
		element.MinVersion = int(rawElement["minVersion"].(float64))
		element.Path = rawElement["path"].(string)
		requestFormat := rawElement["requestFormat"]
		if requestFormat != nil {
			element.RequestFormat = requestFormat.(string)
		}
		elements = append(elements, element)
	}

	return elements, nil

}

func (e *APIInfoElement) String() string {
	return fmt.Sprintf("API:%s maxVersion:%d minVersion:%d path:%s requestFormat:%s", e.API, e.MaxVersion, e.MinVersion, e.Path, e.RequestFormat)
}
