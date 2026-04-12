package version

import (
	"fmt"
	"runtime"
	"strings"
)

const Name = "tun2socks"

var (
	
	
	Version string

	
	
	GitCommit string
)

func String() string {
	return fmt.Sprintf("%s-%s", Name, strings.TrimPrefix(Version, "v"))
}

func BuildString() string {
	return fmt.Sprintf("%s/%s, %s, %s", runtime.GOOS, runtime.GOARCH, runtime.Version(), GitCommit)
}
