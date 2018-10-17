// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

type disk struct {
	data     *synology.Disk
	status   string
	exitcode int
}

type disks []*disk

func CheckDisk(args Args, metrics *Metrics, storageObj *synology.StorageObject) {

	//Required fields
	service := "Disks"
	exitcode := CRITICAL
	message := ""
	perfdata := ""

	// Disk slice
	var disks disks

	// Get synology disks
	synoDisks := storageObj.Disks

	// Fill disk slice
	for i := 0; i < len(synoDisks); i++ {
		nd := new(disk)
		nd.data = &synoDisks[i]
		disks = append(disks, nd)
	}

	// Process disks
	for _, d := range disks {
		setStatusAndExitcode(*d, args.TempWarn, args.TempCrit)
	}

	// Aggregations
	minTemp := disks.minTemp()
	maxTemp := disks.maxTemp()
	avgTemp := disks.avgTemp()

	// Counters
	countDisks := len(disks)
	countOk, countWarning, countCritical, countUnknown := disks.exitcodes()

	// Set overall statuscode
	exitcode = OK
	for _, d := range disks {
		exitcode = maxExitcode(exitcode, d.exitcode)
	}

	switch exitcode {
	case OK:
		message = fmt.Sprintf("Total:%d All disks Ok (Temperature Min=%d\u00b0C Avg=%d\u00b0C Max.%d\u00b0C)", countDisks, minTemp, avgTemp, maxTemp)
	default:
		message = fmt.Sprintf("Total:%d Critical:%d Warning:%d Unknown:%d Ok:%d (Temperature Min=%d\u00b0C Avg=%d\u00b0C Max.%d\u00b0C)", countDisks, countCritical, countWarning, countUnknown, countOk, minTemp, avgTemp, maxTemp)
	}

	perfdata = fmt.Sprintf("Disks_Total=%d Disks_OK=%d Disks_WARNING=%d Disks_CRITICAL=%d Disks_UNKNOWN=%d Temp_Min=%d Temp_Avg=%d Temp_Max=%d", countDisks, countOk, countWarning, countCritical, countUnknown, minTemp, avgTemp, maxTemp)

	// Done. Write the check result
	Write(args, service, exitcode, message, perfdata, metrics)

	// If diskCheck set, create check for each disk
	if args.DiskChecks == true {
		for _, d := range disks {
			var s *synology.Disk
			s = d.data
			// Set servicename
			service = s.Name

			// Set exitcode
			exitcode = d.exitcode

			// Set message
			message = fmt.Sprintf("Type:%s Vendor:%s Model:%s Serial:%s Status:%s Temperature:%d\u00b0C", s.DiskType, s.Vendor, s.Model, s.Serial, d.status, s.Temp)

			// Set perfdata
			perfdata = fmt.Sprintf("Disk_Temperature=%d;%d;%d", s.Temp, args.TempWarn, args.TempCrit)

			// Done. Write the check result
			Write(args, service, exitcode, message, perfdata, metrics)
		}
	}

	return
}

func setStatusAndExitcode(d disk, tempWarn int, tempCrit int) {
	// Status
	status := d.data.Status
	temp := d.data.Temp
	switch status {
	case "normal":
		d.status = status
		d.exitcode = OK
	default:
		d.status = status
		d.exitcode = CRITICAL
	}
	// Temp
	switch {
	case temp < tempWarn:
		d.exitcode = maxExitcode(d.exitcode, OK)
	case temp < tempCrit:
		d.status += " - Temperature-Warning"
		d.exitcode = maxExitcode(d.exitcode, WARNING)
	default:
		d.status += " - Overheating"
		d.exitcode = CRITICAL
	}
}

func (disks *disks) exitcodes() (ok int, warning int, critical int, unknown int) {
	for _, d := range *disks {
		switch d.exitcode {
		case OK:
			ok++
		case WARNING:
			warning++
		case CRITICAL:
			critical++
		case UNKNOWN:
			unknown++
		}
	}
	return ok, warning, critical, unknown
}

/*
func countStatuscode (disks []disk, statuscode int) (c int){
	for i := 0; i < len(disks); i++ {
		if disks[i].Status == statuscode {
			c += 1
		}
	}
	return c
}

func boolToInt(b bool) int {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

func getExitcode (disks []disk) (c int){
	switch {
	case countExitcode(disks, CRITICAL) > 0:
		c = CRITICAL
	case countExitcode(disks, WARNING) > 0:
		c = WARNING
	case countExitcode(disks, UNKNOWN) > 0:
		c = UNKNOWN
	default:
		c = OK
	}
	return c
}

*/

func (disks *disks) minTemp() (t int) {
	t = (*disks)[0].data.Temp
	for i := 0; i < len(*disks); i++ {
		if (*disks)[0].data.Temp < t {
			t = (*disks)[0].data.Temp
		}
	}
	return t
}

func (disks *disks) maxTemp() (t int) {
	t = 0
	for i := 0; i < len(*disks); i++ {
		if (*disks)[0].data.Temp > t {
			t = (*disks)[0].data.Temp
		}
	}
	return t
}

func (disks *disks) avgTemp() (t int) {
	t = 0
	for i := 0; i < len(*disks); i++ {
		t += (*disks)[0].data.Temp
	}
	return int(t / len(*disks))
}
