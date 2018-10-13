// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"
	"strconv"

	"github.com/BFLB/check_synology/synology"
)

func CheckTemperature(args Args, metrics *Metrics, dsmInfo *synology.DSMInfo) {
	// Required fields
	service := "Temperature"
	exitcode := CRITICAL
	message := ""
	perfdata := ""

	// Set Data
	temp := dsmInfo.Temperature
	tempWarning := dsmInfo.TemperatureWarn

	// Set message and perfdata
	message = fmt.Sprintf("%d\u00b0C (Warning:%s)", temp, strconv.FormatBool(tempWarning))
	perfdata = strconv.Itoa(temp)
	perfdata = fmt.Sprintf("Temperature=%d", temp)

	// Set exitcode
	switch tempWarning {
	case false:
		exitcode = OK
	default:
		exitcode = CRITICAL
	}
	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	return

}
