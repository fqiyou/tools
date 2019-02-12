package system

import (
	"github.com/fqiyou/tools/foo/util"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestCollect(t *testing.T)  {
	util.Log.SetLevel(logrus.InfoLevel)
	a := RaidDisk{}
	a.Collect()
	util.JsonPrint(a)

}
