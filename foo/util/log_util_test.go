package util

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRun(t *testing.T) {
	Log.SetLevel(logrus.InfoLevel)
	//ConfigLocalFilesystemLogger("/Users/chaoyang/", "test.log", time.Hour*24*7, time.Hour*1,"%Y%m%d-%H")
	Log.Error("111")
	//Log.Info("222")
	//Log.Info("333")
	//Log.Error("444")

}