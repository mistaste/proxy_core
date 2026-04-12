package version

import (
	"runtime/debug"
)


func Info() []*debug.Module {
	bi, _ := debug.ReadBuildInfo()
	return bi.Deps
}
