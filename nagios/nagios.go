// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"
	"os"
	"time"
)

type Args struct {
	Hostname     string
	Commandfile  string
	UptimeWarn   int
	UptimeCrit   int
	TempWarn     int
	TempCrit     int
	DiskChecks   bool
	PoolWarn     int
	PoolCrit     int
	PoolFailCrit int
	HostPrimary  string
	HostSecondary string
	CPUwarn      int
	CPUcrit      int
	MemWarn      int
	MemCrit      int
}

type Metrics struct {
	Checks      int
	TimeToPrint time.Duration
}

// Nagios exit codes
const OK = 0
const WARNING = 1
const CRITICAL = 2
const UNKNOWN = 3

func Write(a Args, service string, exitcode int, message string, perfdata string, m *Metrics) {
	start := time.Now()
	// Parameters
	command := "PROCESS_SERVICE_CHECK_RESULT"
	timestamp := time.Now().Unix()
	var result string
	// Create the messaga
	if perfdata != "" {
		//result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s: %s | %s", timestamp, command, a.Hostname, service, exitcode, NagiState(exitcode), message, perfdata)
		result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s | %s", timestamp, command, a.Hostname, service, exitcode, message, perfdata)
	} else {
		//result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s", timestamp, command, a.Hostname, service, exitcode, NagiState(exitcode), message)
		result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s", timestamp, command, a.Hostname, service, exitcode, message)
	}

	// TODO Errorhandling?

	if a.Commandfile == "stdout" {
		fmt.Println(result)

	} else {
		f, err := os.OpenFile(a.Commandfile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(f, "%s\n", result)
		f.Close()
	}
	//  Write successful, finally update metrics
	m.Checks += 1
	m.TimeToPrint += time.Now().Sub(start)
}

func NagiState(exitcode int) (state string) {
	switch exitcode {
	case OK:
		return "OK"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func maxExitcode(e1 int, e2 int) int {
	switch e1 {
	case OK:
		return e2
	case WARNING:
		switch e2 {
		case OK:
			return WARNING
		case WARNING:
			return WARNING
		case CRITICAL:
			return CRITICAL
		case UNKNOWN:
			return WARNING
		default:
			return WARNING
		}
	case CRITICAL:
		return CRITICAL
	case UNKNOWN:
		switch e2 {
		case OK:
			return UNKNOWN
		case WARNING:
			return WARNING
		case CRITICAL:
			return CRITICAL
		case UNKNOWN:
			return UNKNOWN
		default:
			return UNKNOWN
		}
	default:
		return UNKNOWN
	}
}
