// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"
	"time"

	"github.com/BFLB/check_synology/synology"
)

func CheckUptime(args Args, metrics *Metrics, dsmInfo *synology.DSMInfo) {

	// Required fields
	service := "Uptime"
	exitcode := OK
	message := ""
	perfdata := ""

	// Get Data
	uptimeSeconds := dsmInfo.Uptime

	// Set message
	message = fmt.Sprintf("%s", time.Duration(uptimeSeconds)*time.Second)

	// Set perfdata
	perfdata = fmt.Sprintf("Uptime=%d", uptimeSeconds)

	// Set exitcode
	switch {
	case uptimeSeconds < args.UptimeCrit:
		exitcode = CRITICAL
	case uptimeSeconds < args.UptimeWarn:
		exitcode = WARNING
	default:
		exitcode = OK
	}

	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	return

}
