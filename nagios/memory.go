// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

func CheckMemory(args Args, metrics *Metrics, util *synology.SystemUtilization) {

	// Required fields
	service := "Memory"
	exitcode := OK
	message := ""
	perfdata := ""

	// Get Data
	availReal := float64(util.Memory.AvailReal) / 1024 / 1024
	availSwap := float64(util.Memory.AvailSwap) / 1024 / 1024
	buffer    := float64(util.Memory.Buffer) / 1024 / 1024
	cached    := float64(util.Memory.Cached) / 1024 / 1024
	memorySize := float64(util.Memory.MemorySize) / 1024 / 1024
	realUsage := util.Memory.RealUsage
	siDisk    := float64(util.Memory.SiDisk) / 1024 / 1024
	soDisk    := float64(util.Memory.SoDisk) / 1024 / 1024
	swapUsage := util.Memory.SwapUsage
	totalReal := float64(util.Memory.TotalReal) / 1024 / 1024
	totalSwap := float64(util.Memory.TotalSwap) /1024 /1024

	message = fmt.Sprintf("Used %d%% Total:%.2fG", realUsage, memorySize)

	perfdata = fmt.Sprintf("Used=%d%%;%d;%d Avail_Real=%.2fG AvailSwap=%.2fG Buffer=%.2fG Cached=%.2fG MemorySize=%.2fG SiDisk=%.2fG SoDisk=%.2fG SwapUsage=%d%% TotalReal=%.2fG TotalSwap=%.2fG", realUsage, args.MemWarn, args.MemCrit, availReal, availSwap, buffer, cached, memorySize, siDisk, soDisk, swapUsage, totalReal, totalSwap )

	// Set exitcode
	switch {
	case realUsage >= args.MemCrit:
		exitcode = CRITICAL
	case realUsage >= args.MemWarn:
		exitcode = WARNING
	default:
		exitcode = OK
	}

	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	return

}

