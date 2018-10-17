// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BFLB/check_synology/synology"
)

type pool struct {
	data        *synology.StoragePool
	name        string
	percentUsed int
	sizeFree    int
	sizeUsed    int
	sizeTotal   int
	status      string
	exitcode    int
}

type pools []*pool

func CheckStoragePool(args Args, metrics *Metrics, storageObj *synology.StorageObject) {

	//Required fields
	service := "Storagepool"
	exitcode := CRITICAL
	message := ""
	perfdata := ""

	// Pool slice
	var pools pools

	// Get synology pools
	synoPools := storageObj.StoragePools

	// Fill pool slice
	for i := 0; i < len(synoPools); i++ {
		np := new(pool)
		np.data = &synoPools[i]
		pools = append(pools, np)
	}

	// Setup pools
	for _, p := range pools {
		p.name = strings.Replace(string(p.data.ID), "reuse_", "Storagepool ", 1)
		used, _ := strconv.Atoi(p.data.Size.Used) // TODO Errorhandling
		p.sizeUsed = used / 1024 / 1024 / 1024
		total, _ := strconv.Atoi(p.data.Size.Total)
		p.sizeTotal = total / 1024 / 1024 / 1024 // TODO Errorhandling
		p.sizeFree = p.sizeTotal - p.sizeUsed
		p.percentUsed = p.sizeUsed / p.sizeTotal * 100
	}

	// Process
	for _, p := range pools {
		service = p.name
		message = fmt.Sprintf("Status:%s Type:%s Size-Used:%dGB(%d%%) Size-Free:%dGB Size-Total:%dGB Disk-failures:%d", p.status, p.data.DeviceType, p.sizeUsed, p.percentUsed, p.sizeFree, p.sizeTotal, p.data.DiskFailureNumber)
		perfdata = fmt.Sprintf("Size_Used_Percent=%d%%;%d;%d Size_Used=%dG Size_Free=%dG Size_Total=%dG Disk_Failures=%d;;%d", p.percentUsed, args.PoolWarn, args.PoolCrit, p.sizeUsed, p.sizeFree, p.sizeTotal, p.data.DiskFailureNumber, args.PoolFailCrit)

		// Status
		status := p.data.Status
		failed := p.data.DiskFailureNumber
		usage := p.percentUsed
		switch {
		case failed >= args.PoolFailCrit:
			p.status = status
			p.exitcode = OK
		case failed > 0:
			p.status = status
			p.exitcode = WARNING
		default:
			p.status = status
			p.exitcode = OK
		}
		// Space
		switch {
		case usage >= args.PoolCrit:
			p.status += " - Sice Critical"
			p.exitcode = CRITICAL
		case usage >= args.PoolWarn:
			p.status += " - Size Warning"
			p.exitcode = maxExitcode(p.exitcode, WARNING)
		default:
			p.exitcode = maxExitcode(p.exitcode, OK)
		}

		// Done. Write the check result
		Write(args, service, exitcode, message, perfdata, metrics)
	}

	return
}
