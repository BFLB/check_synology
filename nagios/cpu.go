// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

func CheckCPU(args Args, metrics *Metrics, util *synology.SystemUtilization) {

	// Required fields
	service := "CPU"
	exitcode := OK
	message := ""
	perfdata := ""

	// Get Data
	fifteenMinLoad := float32(util.CPU.FifteenMinLoad) / 100 // Weird but true
	oneMinLoad     := float32(util.CPU.OneMinLoad) / 100  // Weird but true
	fiveMinLoad    := float32(util.CPU.FiveMinLoad) / 100 // Weird but true
	otherLoad      := util.CPU.OtherLoad
	systemLoad     := util.CPU.SystemLoad
	userLoad       := util.CPU.UserLoad
	currentLoad    := systemLoad + userLoad + otherLoad

	message = fmt.Sprintf("%d%% (User:%d%% System:%d%% Other:%d%%) Load1min:%.2f Load5min:%.2f Load15min:%.2f", currentLoad, userLoad, systemLoad, otherLoad, oneMinLoad, fiveMinLoad, fifteenMinLoad  )

	perfdata = fmt.Sprintf("CPU_Total=%d%%,%d,%d CPU_User=%d%% CPU_System=%d%% CPU_Other=%d%% Load_1min=%.2f Load_5min=%.2f Load_15min:%.2f", currentLoad, args.CPUwarn, args.CPUcrit, userLoad, systemLoad, otherLoad, oneMinLoad, fiveMinLoad, fifteenMinLoad  )

	// Set exitcode
	switch {
	case currentLoad >= args.CPUcrit:
		exitcode = CRITICAL
	case currentLoad >= args.CPUwarn:
		exitcode = WARNING
	default:
		exitcode = OK
	}

	//fmt.Printf("Device: %s\n", device)

	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	return

}

