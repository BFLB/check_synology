// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

func CheckSystemStatus(args Args, metrics *Metrics, sysStat *synology.SystemStatus) {

	// Required fields
	service := "SystemStatus"
	exitcode := OK
	message := ""
	perfdata := ""

	// Get Data
	crashed := sysStat.IsSystemCrashed
	upgradeReady := sysStat.UpgradeReady

	// Set message, exitcode and perfdata
	switch crashed {
	case true:
		message = fmt.Sprintf("Status:Crashed")
		exitcode = CRITICAL
	case false:
		message = fmt.Sprintf("Status:OK")
		exitcode = OK
	}
	switch upgradeReady {
	case true:
		message += " - Upgrade Ready"
		switch exitcode {
		case CRITICAL:
			// NOOP
		default:
			exitcode = WARNING
		}
	default:
		// NOOP

	}

	perfdata = fmt.Sprintf("Crashed=%d UpgradeAvailable=%d", boolToInt(crashed), boolToInt(upgradeReady))

	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	return

}

func boolToInt(b bool) int {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}
