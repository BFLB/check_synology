// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

// Example command add-contact
// Shows 3CX system information
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BFLB/check_synology/nagios"
	synology "github.com/BFLB/check_synology/synology"
)

// Version
const _version = "0.2"

// Comman-line Arguments
var (
	version       = flag.Bool("V", false, "Show version")
	host          = flag.String("H", "", "Synology Address")
	user          = flag.String("u", "admin", "Username")
	pass          = flag.String("a", "admin", "Password")
	port          = flag.String("p", "5001", "Port")
	timeout       = flag.Int("t", 10, "Timeout")
	commandfile   = flag.String("cmd", "stdout", "Nagios command file")
	diskChecks    = flag.Bool("d", true, "Individual check for each disk")
	extCheck      = flag.Bool("extCheck", false, "Checks for extension-units")
	hostPrimary   = flag.String("prim", "", "SN of primary host")
	hostSecondary = flag.String("sec", "", "SN of secondary host")
	uptimeWarn    = flag.Int("wUp", 86400, "Uptime Warning (s)")
	uptimeCrit    = flag.Int("cUp", 0, "Uptime Warning (s)")
	tempWarn      = flag.Int("wTemp", 40, "Temperature Warning")
	tempCrit      = flag.Int("cTemp", 50, "Temperature Critical")
	poolWarn      = flag.Int("wPool", 80, "Used % warning")
	poolCrit      = flag.Int("cPool", 90, "Used % critical")
	poolFailCrit  = flag.Int("cPoolFail", 1, "Number of failed disks per RAID for state critical")
)

func main() {

	// Parte command-line args
	flag.Parse()

	// Print version information
	if *version == true {
		fmt.Printf("Version: %s", _version)
		os.Exit(0)
	}

	exitcode := nagios.CRITICAL

	tExecCrit := *timeout
	tExecWarn := int((tExecCrit / 10) * 8)

	// Timestamp to calculate execution time
	timestampStart := time.Now()

	// Login to 3CX pbx
	api, err := synology.Login(*user, *pass, *host, *port)
	if err != nil {
		log.Fatal("Login returned error: ", err)
	}
	defer api.Logout(*user, *pass)

	timestampLogin := time.Now()

	// Fetch Synology Information
	var dsmInfo *synology.DSMInfo
	dsmInfo, err = api.DSMInfo()
	if err != nil {
		exitcode = nagios.CRITICAL
		fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
		os.Exit(exitcode)
	}
	var systemStatus *synology.SystemStatus
	systemStatus, err = api.SystemStatus()
	if err != nil {
		exitcode = nagios.CRITICAL
		fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
		os.Exit(exitcode)
	}
	/*
		var systemUtilization *synology.SystemUtilization
		systemUtilization, err = api.SystemUtilization()
		if err != nil {
			exitcode = nagios.CRITICAL
			fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
			os.Exit(exitcode)
		}
	*/

	var storage *synology.StorageObject
	storage, err = api.Storage()
	if err != nil {
		exitcode = nagios.CRITICAL
		fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
		os.Exit(exitcode)
	}
	//fmt.Printf("Storage:\n\n\n%v\n\n\n", storage) // FIXME Remove

/*
		var apiInfo []synology.APIInfoElement
		apiInfo, err = api.APIInfo()
		if err != nil {
			exitcode = nagios.CRITICAL
			fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
			os.Exit(exitcode)
		}
		fmt.Println("\nAPI-INFO")
		for _, e := range apiInfo {
			fmt.Printf("    %s\n", e.String())
		}
*/	
	var dualEnclosure *synology.DualEnclosure
	dualEnclosure, err = api.DualEnclosure()
	if err != nil {
		//exitcode = nagios.CRITICAL
		//fmt.Printf("%s: Plugin version: %s - %s\n", nagios.NagiState(exitcode), _version, err.Error())
		//os.Exit(exitcode)
		dualEnclosure = new(synology.DualEnclosure) // FIXME
	}

	timestampFetch := time.Now()

	// Setup Nagios
	var nagiArgs nagios.Args
	nagiArgs.Hostname = *host
	nagiArgs.Commandfile = *commandfile
	nagiArgs.UptimeWarn = *uptimeWarn
	nagiArgs.UptimeCrit = *uptimeCrit
	nagiArgs.TempWarn = *tempWarn
	nagiArgs.TempCrit = *tempCrit
	nagiArgs.DiskChecks = *diskChecks
	nagiArgs.PoolWarn = *poolWarn
	nagiArgs.PoolCrit = *poolCrit
	nagiArgs.PoolFailCrit = *poolFailCrit
	nagiArgs.HostPrimary = *hostPrimary
	nagiArgs.HostSecondary = *hostSecondary

	var nagiMetrics nagios.Metrics

	// Do the checks
	nagios.CheckModel(nagiArgs, &nagiMetrics, dsmInfo)
	nagios.CheckSystemStatus(nagiArgs, &nagiMetrics, systemStatus)
	nagios.CheckTemperature(nagiArgs, &nagiMetrics, dsmInfo)
	nagios.CheckUptime(nagiArgs, &nagiMetrics, dsmInfo)
	nagios.CheckDisk(nagiArgs, &nagiMetrics, storage)
	nagios.CheckStoragePool(nagiArgs, &nagiMetrics, storage)
	nagios.CheckEnclosure(nagiArgs, &nagiMetrics, dualEnclosure)

	timestampProcess := time.Now()

	tLogin := timestampLogin.Sub(timestampStart).Seconds()
	tFetch := timestampFetch.Sub(timestampLogin).Seconds()
	tWrite := nagiMetrics.TimeToPrint.Seconds() // FIXME TimeToWrite
	tProcess := (timestampProcess.Sub(timestampFetch) - nagiMetrics.TimeToPrint).Seconds()
	tExec := timestampProcess.Sub(timestampStart).Seconds()

	// Prepare exit information
	// Set exitcode
	if tExec > float64(*timeout) {
		exitcode = nagios.CRITICAL

	} else if tExec >= float64(*timeout) {
		exitcode = nagios.WARNING
	} else {
		exitcode = nagios.OK
	}

	message := fmt.Sprintf("%d passive-check(s) generated in %.3f seconds (t_conn:%.3fs t_fetch:%.3fs, t_proc:%.3fs, t_print:%.3fs)", nagiMetrics.Checks, tExec, tLogin, tFetch, tProcess, tWrite)

	perfdata := fmt.Sprintf("ExecTime=%3.3fs;%d;%d t_conn=%3.3fs t_load=%3.3fs t_proc=%3.3fs t_print=%3.3fs StatusCode=%d ChecksCreated=%d", tExec, tExecWarn, tExecCrit, tLogin, tFetch, tProcess, tWrite, exitcode, nagiMetrics.Checks)

	// TODO Make exit function in utils
	// Print exit information
	fmt.Printf("%s: Plugin version: %s - %s | %s\n", nagios.NagiState(exitcode), _version, message, perfdata)

	// Done. Exit with exitcode
	os.Exit(exitcode)

}
