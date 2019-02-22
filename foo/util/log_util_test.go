package util

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRun(t *testing.T) {
	// netstat -aulntp | grep rsyslog
	// /etc/rsyslog.conf

	app_name := "dev_golang_app"

	Log.SetLevel(logrus.InfoLevel)

	//AddHookLocalFileLogger("/Users/chaoyang/tmp/", app_name+".log", time.Hour*24*7, time.Hour*1,"%Y%m%d-%H")
	//AddHookSyslog("udp","spark003:514",syslog.LOG_LOCAL4,app_name)


	Log.Error("111")
	Log.Info("222")
	Log.Info("333")
	Log.Error(app_name)

}