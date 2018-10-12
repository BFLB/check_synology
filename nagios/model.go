// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package nagios

import (
	"fmt"

	"github.com/BFLB/check_synology/synology"
)

func CheckModel(args Args, metrics *Metrics, dsmInfo *synology.DSMInfo) {
	service := "Model"
	exitcode := OK
	perfdata := ""

	var message string
	message = fmt.Sprintf("%s (S/N:%s)", dsmInfo.Model, dsmInfo.Serial)
	Write(args, service, exitcode, message, perfdata, metrics)

	return
}
