package util

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	Log.Info(NewUUID())
}