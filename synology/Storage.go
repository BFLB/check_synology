// Copyright (c) 2018 Bernhard Fluehmann. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.
//
// Library to access Synology DSM API
//

package synology

import (
	"encoding/json"
	"errors"
	"net/url"
	//"net/http/httputil"
)

type storageResponse struct {
	Data    StorageObject `json:"data"`
	Error   struct {
		Code int `json:"code"`
	} `json:"error, omitempty"`
	Success bool         `json:"success"`
}

type StorageObject struct {
	AHAInfo struct {
		EncCnt     int `json:"enc_cnt"`
		Enclosures []struct {
			Disks       []int  `json:"disks"`
			Fans        []int  `json:"fans"`
			Links       []int  `json:"links"`
			MaxDisk     int    `json:"max_disk"`
			Model       string `json:"model"`
			ModelID     int    `json:"model_id"`
			Powers      []int  `json:"powers"`
			Temperature int    `json:"temperature"`
		} `json:"enclosures"`
		HostCnt int `json:"host_cnt"`
		Hosts   []struct {
			HostType  int    `json:"host_type"`
			Links     []int  `json:"links"`
			Modelname string `json:"modelname"`
		} `json:"hosts"`
		LinkStatus int `json:"link_status"`
	} `json:"AHAInfo"`
	Disks []Disk `json:"disks"`
	Enclosures []struct {
		Disks []struct {
			ID string `json:"id"`
		} `json:"disks"`
		ID       int  `json:"id"`
		Internal bool `json:"internal"`
		Ports    []struct {
			LinkID      string `json:"linkId"`
			LinkPortNum int    `json:"linkPortNum"`
			PortNum     int    `json:"portNum"`
			Valid       bool   `json:"valid"`
		} `json:"ports"`
	} `json:"enclosures"`
	Env struct {
		Batchtask struct {
			MaxTask    int `json:"max_task"`
			RemainTask int `json:"remain_task"`
		} `json:"batchtask"`
		BayNumber     string `json:"bay_number"`
		DataScrubbing struct {
			ScheEnabled string `json:"sche_enabled"`
			ScheStatus  string `json:"sche_status"`
		} `json:"data_scrubbing"`
		Ebox               []interface{} `json:"ebox"`
		FsActing           bool          `json:"fs_acting"`
		IsSyncSysPartition bool          `json:"isSyncSysPartition"`
		IsSpaceActioning   bool          `json:"is_space_actioning"`
		Isns               struct {
			Address string `json:"address"`
			Enabled bool   `json:"enabled"`
		} `json:"isns"`
		IsnsServer            string `json:"isns_server"`
		MaxFsBytes            string `json:"max_fs_bytes"`
		MaxFsBytesHighEnd     string `json:"max_fs_bytes_high_end"`
		ModelName             string `json:"model_name"`
		RAMEnoughForFsHighEnd bool   `json:"ram_enough_for_fs_high_end"`
		RAMSize               int    `json:"ram_size"`
		RAMSizeRequired       int    `json:"ram_size_required"`
		Showpooltab           bool   `json:"showpooltab"`
		Status                struct {
			SystemCrashed    bool `json:"system_crashed"`
			SystemNeedRepair bool `json:"system_need_repair"`
		} `json:"status"`
		Support struct {
			Ebox      bool `json:"ebox"`
			RaidCross bool `json:"raid_cross"`
			Sysdef    bool `json:"sysdef"`
		} `json:"support"`
		SupportFitFsLimit  bool    `json:"support_fit_fs_limit"`
		UniqueKey          string  `json:"unique_key"`
		VolumeFullCritical float64 `json:"volume_full_critical"`
		VolumeFullWarning  float64 `json:"volume_full_warning"`
	} `json:"env"`
	HotSpareConf struct {
		CrossRepair   bool          `json:"cross_repair"`
		DisableRepair []interface{} `json:"disable_repair"`
	} `json:"hotSpareConf"`
	HotSpares []interface{} `json:"hotSpares"`
	IscsiLuns []struct {
		CanDo struct {
			DataScrubbing bool `json:"data_scrubbing"`
			Delete        bool `json:"delete"`
			ExpandByDisk  int  `json:"expand_by_disk"`
			RaidCross     bool `json:"raid_cross"`
		} `json:"can_do"`
		ID          string `json:"id"`
		IsActioning bool   `json:"is_actioning"`
		IscsiLun    struct {
			BlkNum        string `json:"blkNum"`
			DeviceType    string `json:"device_type"`
			ExtentBased   bool   `json:"extent_based"`
			ExtentSize    string `json:"extent_size"`
			IsBlun        bool   `json:"is_blun"`
			Lid           int    `json:"lid"`
			Location      string `json:"location"`
			MappedTargets []int  `json:"mapped_targets"`
			Name          string `json:"name"`
			RestoredTime  string `json:"restored_time"`
			Rootpath      string `json:"rootpath"`
			Size          string `json:"size"`
			ThinProvision bool   `json:"thin_provision"`
			UsedBy        string `json:"used_by"`
			UUID          string `json:"uuid"`
		} `json:"iscsi_lun"`
		NumID    int `json:"num_id"`
		Progress struct {
			Percent string `json:"percent"`
			Step    string `json:"step"`
		} `json:"progress"`
		Status string `json:"status"`
	} `json:"iscsiLuns"`
	IscsiTargets []struct {
		Auth struct {
			MutualUsername string `json:"mutual_username"`
			Type           string `json:"type"`
			Username       string `json:"username"`
		} `json:"auth"`
		DataChksum              bool   `json:"data_chksum"`
		Enabled                 bool   `json:"enabled"`
		HdrChksum               bool   `json:"hdr_chksum"`
		Iqn                     string `json:"iqn"`
		MappedLogicalUnitNumber []int  `json:"mapped_logical_unit_number"`
		MappedLuns              []int  `json:"mapped_luns"`
		Masking                 []struct {
			Iqn        string `json:"iqn"`
			Permission string `json:"permission"`
		} `json:"masking"`
		MultiSessions bool   `json:"multi_sessions"`
		Name          string `json:"name"`
		NumID         int    `json:"num_id"`
		RecvSegBytes  int    `json:"recv_seg_bytes"`
		Remote        []struct {
			IP  string `json:"ip"`
			Iqn string `json:"iqn"`
		} `json:"remote"`
		SendSegBytes int    `json:"send_seg_bytes"`
		Status       string `json:"status"`
		Tid          int    `json:"tid"`
	} `json:"iscsiTargets"`
	Ports        []interface{} `json:"ports"`
	SsdCaches    []interface{} `json:"ssdCaches"`
	StoragePools []struct {
		CacheStatus string `json:"cacheStatus"`
		CanDo       struct {
			DataScrubbing bool `json:"data_scrubbing"`
			Delete        bool `json:"delete"`
			ExpandByDisk  int  `json:"expand_by_disk"`
			RaidCross     bool `json:"raid_cross"`
		} `json:"can_do"`
		Container         string   `json:"container"`
		DeployPath        string   `json:"deploy_path"`
		Desc              string   `json:"desc"`
		DeviceType        string   `json:"device_type"`
		DiskFailureNumber int      `json:"disk_failure_number"`
		Disks             []string `json:"disks"`
		DriveType         int      `json:"drive_type"`
		ID                string   `json:"id"`
		IsActioning       bool     `json:"is_actioning"`
		IsScheduled       bool     `json:"is_scheduled"`
		IsWritable        bool     `json:"is_writable"`
		LastDoneTime      int      `json:"last_done_time"`
		LimitedDiskNumber int      `json:"limited_disk_number"`
		MaximalDiskSize   string   `json:"maximal_disk_size"`
		MinimalDiskSize   string   `json:"minimal_disk_size"`
		NextScheduleTime  int      `json:"next_schedule_time"`
		NumID             int      `json:"num_id"`
		PoolPath          string   `json:"pool_path"`
		Progress          struct {
			Percent string `json:"percent"`
			Step    string `json:"step"`
		} `json:"progress"`
		RaidType string `json:"raidType"`
		Raids    []struct {
			DesignedDiskCount int `json:"designedDiskCount"`
			Devices           []struct {
				ID     string `json:"id"`
				Slot   int    `json:"slot"`
				Status string `json:"status"`
			} `json:"devices"`
			HasParity      bool          `json:"hasParity"`
			MinDevSize     string        `json:"minDevSize"`
			NormalDevCount int           `json:"normalDevCount"`
			RaidPath       string        `json:"raidPath"`
			RaidStatus     int           `json:"raidStatus"`
			Spares         []interface{} `json:"spares"`
		} `json:"raids"`
		ScrubbingStatus string `json:"scrubbingStatus"`
		Size            struct {
			Total string `json:"total"`
			Used  string `json:"used"`
		} `json:"size"`
		SpacePath string `json:"space_path"`
		SsdTrim   struct {
			Support string `json:"support"`
		} `json:"ssd_trim"`
		Status      string        `json:"status"`
		Suggestions []interface{} `json:"suggestions"`
		Timebackup  bool          `json:"timebackup"`
		VspaceCanDo struct {
			Drbd struct {
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"drbd"`
			Flashcache struct {
				Apply struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"apply"`
				Remove struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"remove"`
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"flashcache"`
			Snapshot struct {
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"snapshot"`
		} `json:"vspace_can_do"`
	} `json:"storagePools"`
	Volumes []struct {
		AtimeChecked bool   `json:"atime_checked"`
		AtimeOpt     string `json:"atime_opt"`
		CacheStatus  string `json:"cacheStatus"`
		CanDo        struct {
			DataScrubbing bool `json:"data_scrubbing"`
			Delete        bool `json:"delete"`
			ExpandByDisk  int  `json:"expand_by_disk"`
			RaidCross     bool `json:"raid_cross"`
		} `json:"can_do"`
		Container         string        `json:"container"`
		DeployPath        string        `json:"deploy_path"`
		Desc              string        `json:"desc"`
		DeviceType        string        `json:"device_type"`
		DiskFailureNumber int           `json:"disk_failure_number"`
		Disks             []interface{} `json:"disks"`
		DriveType         int           `json:"drive_type"`
		EppoolUsed        string        `json:"eppool_used"`
		ExistAliveVdsm    bool          `json:"exist_alive_vdsm"`
		FsType            string        `json:"fs_type"`
		ID                string        `json:"id"`
		IsActing          bool          `json:"is_acting"`
		IsActioning       bool          `json:"is_actioning"`
		IsInodeFull       bool          `json:"is_inode_full"`
		IsScheduled       bool          `json:"is_scheduled"`
		IsWritable        bool          `json:"is_writable"`
		LastDoneTime      int           `json:"last_done_time"`
		LimitedDiskNumber int           `json:"limited_disk_number"`
		MaxFsSize         string        `json:"max_fs_size"`
		NextScheduleTime  int           `json:"next_schedule_time"`
		NumID             int           `json:"num_id"`
		PoolPath          string        `json:"pool_path"`
		Progress          struct {
			Percent string `json:"percent"`
			Step    string `json:"step"`
		} `json:"progress"`
		RaidType        string `json:"raidType"`
		ScrubbingStatus string `json:"scrubbingStatus"`
		Size            struct {
			FreeInode   string `json:"free_inode"`
			Total       string `json:"total"`
			TotalDevice string `json:"total_device"`
			TotalInode  string `json:"total_inode"`
			Used        string `json:"used"`
		} `json:"size"`
		SsdTrim struct {
			Support string `json:"support"`
		} `json:"ssd_trim"`
		Status        string        `json:"status"`
		Suggestions   []interface{} `json:"suggestions"`
		Timebackup    bool          `json:"timebackup"`
		UsedByGluster bool          `json:"used_by_gluster"`
		VolPath       string        `json:"vol_path"`
		VspaceCanDo   struct {
			Drbd struct {
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"drbd"`
			Flashcache struct {
				Apply struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"apply"`
				Remove struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"remove"`
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"flashcache"`
			Snapshot struct {
				Resize struct {
					CanDo       bool `json:"can_do"`
					ErrCode     int  `json:"errCode"`
					StopService bool `json:"stopService"`
				} `json:"resize"`
			} `json:"snapshot"`
		} `json:"vspace_can_do"`
	} `json:"volumes"`
}

type Disk struct {
		AdvProgress        string `json:"adv_progress"`
		AdvStatus          string `json:"adv_status"`
		BelowRemainLifeThr bool   `json:"below_remain_life_thr"`
		Container          struct {
			Order int    `json:"order"`
			Str   string `json:"str"`
			Type  string `json:"type"`
		} `json:"container"`
		Device             string `json:"device"`
		DisableSecera      bool   `json:"disable_secera"`
		DiskType           string `json:"diskType"`
		DiskCode           string `json:"disk_code"`
		EraseTime          int    `json:"erase_time"`
		ExceedBadSectorThr bool   `json:"exceed_bad_sector_thr"`
		Firm               string `json:"firm"`
		HasSystem          bool   `json:"has_system"`
		ID                 string `json:"id"`
		IhmTesting         bool   `json:"ihm_testing"`
		Is4Kn              bool   `json:"is4Kn"`
		IsSsd              bool   `json:"isSsd"`
		IsSynoPartition    bool   `json:"isSynoPartition"`
		IsErasing          bool   `json:"is_erasing"`
		LongName           string `json:"longName"`
		Model              string `json:"model"`
		Name               string `json:"name"`
		NumID              int    `json:"num_id"`
		Order              int    `json:"order"`
		OverviewStatus     string `json:"overview_status"`
		PciSlot            int    `json:"pciSlot"`
		PerfTesting        bool   `json:"perf_testing"`
		PortType           string `json:"portType"`
		RemainLife         int    `json:"remain_life"`
		Serial             string `json:"serial"`
		SizeTotal          string `json:"size_total"`
		SmartProgress      string `json:"smart_progress"`
		SmartStatus        string `json:"smart_status"`
		SmartTestLimit     int    `json:"smart_test_limit"`
		SmartTesting       bool   `json:"smart_testing"`
		Status             string `json:"status"`
		Support            bool   `json:"support"`
		Temp               int    `json:"temp"`
		TestingProgress    string `json:"testing_progress"`
		TestingType        string `json:"testing_type"`
		TrayStatus         string `json:"tray_status"`
		Unc                int    `json:"unc"`
		UsedBy             string `json:"used_by"`
		Vendor             string `json:"vendor"`
	}


func (api *Syno) Storage() (*StorageObject, error) {
	// Set URL parameters
	parameters := url.Values{}
	parameters.Add("api", "SYNO.Storage.CGI.Storage")
	parameters.Add("version", "1")
	parameters.Add("method", "load_info")

	response, err := api.Get("entry.cgi", &parameters)
	if err != nil {
		return nil, err
	}

	var payload *storageResponse
	err = json.Unmarshal(response, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == false {
		return nil, errors.New(string(response))
	}

	return &payload.Data, nil

}
