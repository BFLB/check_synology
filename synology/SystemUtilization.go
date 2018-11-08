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

type systemUtilizationResponse struct {
	Data    SystemUtilization `json:"data"`
	Success bool              `json:"success"`
}

type SystemUtilization struct {
	CPU struct {
		FifteenMinLoad int    `json:"15min_load"`
		OneMinLoad     int   `json:"1min_load"`
		FiveMinLoad    int    `json:"5min_load"`
		Device         string `json:"device"`
		OtherLoad      int    `json:"other_load"`
		SystemLoad     int    `json:"system_load"`
		UserLoad       int    `json:"user_load"`
	} `json:"cpu"`
	Disk struct {
		Disk []struct {
			Device      string `json:"device"`
			DisplayName string `json:"display_name"`
			ReadAccess  int    `json:"read_access"`
			ReadByte    int    `json:"read_byte"`
			DiskType    string `json:"type"`
			Utilization int    `json:"utilization"`
			WriteAccess int    `json:"write_access"`
			WriteByte   int    `json:"write_byte"`
		} `json:"disk"`
		Total struct {
			Device      string `json:"device"`
			ReadAccess  int    `json:"read_access"`
			ReadByte    int    `json:"read_byte"`
			Utilization int    `json:"utilization"`
			WriteAccess int    `json:"write_access"`
			WriteByte   int    `json:"write_byte"`
		} `json:"total"`
	} `json:"disk"`
	Lun    []interface{} `json:"lun"`
	Memory struct {
		AvailReal  int    `json:"avail_real"`
		AvailSwap  int    `json:"avail_swap"`
		Buffer     int    `json:"buffer"`
		Cached     int    `json:"cached"`
		Device     string `json:"device"`
		MemorySize int    `json:"memory_size"`
		RealUsage  int    `json:"real_usage"`
		SiDisk     int    `json:"si_disk"`
		SoDisk     int    `json:"so_disk"`
		SwapUsage  int    `json:"swap_usage"`
		TotalReal  int    `json:"total_real"`
		TotalSwap  int    `json:"total_swap"`
	} `json:"memory"`
	Network []struct {
		Device string `json:"device"`
		Rx     int    `json:"rx"`
		Tx     int    `json:"tx"`
	} `json:"network"`
	Space struct {
		Total struct {
			Device      string `json:"device"`
			ReadAccess  int    `json:"read_access"`
			ReadByte    int    `json:"read_byte"`
			Utilization int    `json:"utilization"`
			WriteAccess int    `json:"write_access"`
			WriteByte   int    `json:"write_byte"`
		} `json:"total"`
		Volume []struct {
			Device      string `json:"device"`
			DisplayName string `json:"display_name"`
			ReadAccess  int    `json:"read_access"`
			ReadByte    int    `json:"read_byte"`
			Utilization int    `json:"utilization"`
			WriteAccess int    `json:"write_access"`
			WriteByte   int    `json:"write_byte"`
		} `json:"volume"`
	} `json:"space"`
	Time int `json:"time"`
}

func (api *Syno) SystemUtilization() (*SystemUtilization, error) {

	// Set URL parameters
	//api":"SYNO.Core.System.Utilization","method":"get","version":1,"type":"current","resource":["cpu","memory","network","lun","disk","space"]}]
	parameters := url.Values{}
	parameters.Add("api", "SYNO.Core.System.Utilization")
	parameters.Add("version", "1")
	parameters.Add("method", "get")

	response, err := api.Get("entry.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	var payload *systemUtilizationResponse
	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == false {
		return nil, errors.New(string(response))
	}

	return &payload.Data, nil

}
