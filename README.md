# check_synology

Golang Nagios Plugin for Synology devices. Uses the DSM Web API.

This plugin can be executed by any Nagios compatible Monitoring System as an actice check. During the execution, the various results are sent to the Monitoring System as independent passive checks. At the end, the plugin returns with some information of the overall execution results, e.g. Execution time, Number of checks created etc.

Feedback and contribution welcome!

Checks:
- Model / Serialnumber / Version
- systemStatus
- Temperature
- Uptime
- Disks
- Storagepools (RAIDS)
- Hosts (Controller)
- Enclosures (Extension Unit)
- CPU
- Memory
- Network

