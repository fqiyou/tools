package system

import (
	"github.com/fqiyou/tools/foo/util"
	"os"
	"strconv"
	"strings"
)

type RaidDiskSystemInfo struct {
	RaidName    	string
	RaidLevel		string
	Version 		string
	RaidDevices		uint64
	TotalDevices 	uint64
	ActiveDevices 	uint64
	WorkingDevices 	uint64
	FailedDevices 	uint64
	SpareDevices 	uint64
}

type RaidDisk struct {
	RaidDiskSystemInfo 		map[string]RaidDiskSystemInfo
	HostName 				string
}



func (this *RaidDisk) Collect() error {

	//cmd := "cat /Users/chaoyang/GoProject/src/github.com/fqiyou/tools/foo/system/mdstat |grep md|awk -F \":\" '{print $1}'"
	cmd := "mdstat |grep md|awk -F \":\" '{print $1}'"
	output, err := Exec(cmd)
	if err != nil {
		util.Log.Error(err)
		return err
	}

	this.RaidDiskSystemInfo = map[string]RaidDiskSystemInfo{}

	lines := strings.Split(output, "\n")
	for _, line := range lines{
		Trim(&line)
		if len(line) == 0 {
			continue
		}
		this.RaidDiskSystemInfo[line] = getRaidDiskInfo(line)
	}
	this.HostName, _ = os.Hostname()
	return nil
}




func getRaidDiskInfo(mount string) RaidDiskSystemInfo {
	rdsi := RaidDiskSystemInfo{}
	rdsi.RaidName = mount

	raid_level_cmd := getMdadmCmd("Raid Level",mount)
	rdsi.RaidLevel = ExecOutput(raid_level_cmd)
	version_cmd := getMdadmCmd("Version",mount)
	rdsi.Version = ExecOutput(version_cmd)


	raid_devices := getMdadmCmd("Raid Devices",mount)
	rdsi.RaidDevices,_ = strconv.ParseUint(ExecOutput(raid_devices), 10, 64)

	total_devices := getMdadmCmd("Total Devices",mount)
	rdsi.TotalDevices,_ = strconv.ParseUint(ExecOutput(total_devices), 10, 64)

	active_devices := getMdadmCmd("Active Devices",mount)
	rdsi.ActiveDevices,_ = strconv.ParseUint(ExecOutput(active_devices), 10, 64)

	working_devices := getMdadmCmd("Working Devices",mount)
	rdsi.WorkingDevices,_ = strconv.ParseUint(ExecOutput(working_devices), 10, 64)

	failed_devices := getMdadmCmd("Failed Devices",mount)
	rdsi.FailedDevices,_ = strconv.ParseUint(ExecOutput(failed_devices), 10, 64)

	spare_devices := getMdadmCmd("Spare Devices",mount)
	rdsi.SpareDevices,_ = strconv.ParseUint(ExecOutput(spare_devices), 10, 64)

	return rdsi
}

func getMdadmCmd(args string,mount string) string{
	//cmd := "cat /Users/chaoyang/GoProject/src/github.com/fqiyou/tools/foo/system/mdadm_cmd |grep \"" + args + "\"|awk -F \":\" '{print $2}' |tr -d \" \""

	cmd := "mdadm -D /dev/" + mount + " |grep \"" + args + "\"|awk -F \":\" '{print $2}' |tr -d \" \""

	return cmd
}
