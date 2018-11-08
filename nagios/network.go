// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

func CheckNetwork(args Args, metrics *Metrics, util *synology.SystemUtilization) {

	// Required fields
	service := "Memory"
	exitcode := OK
	message := ""
	perfdata := ""

	// Get Data
	var ports []synology.NetworkPort
	ports = util.Network
	
	// Fill disk slice
	for i := 0; i < len(ports); i++ {
		service = fmt.Sprintf("Network %s", ports[i].Device) 
		rx := float64(ports[i].Rx) / 1024 / 1024
		tx := float64(ports[i].Tx) / 1024 / 1024
		message = fmt.Sprintf("Rx:%.3fGb/s Tx:%.3fGb/s", rx, tx )

		perfdata = fmt.Sprintf("Rx=%.3fG Tx=%.3fG", rx, tx )
		
		// Write the check result
		Write(args, service, exitcode, message, perfdata, metrics)
	}

	return

}

