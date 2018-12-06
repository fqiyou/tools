package util

import (
	"github.com/fqiyou/tools/foo/tools/logs"
	"testing"
)

func TestNewUUID(t *testing.T) {
	log.NewLogger(10000)
	log.SetLevel(log.LevelDebug)
	log.EnableFuncCallDepth(true)

	log.SetLogger("console", "")
	log.Info(NewUUID())

}