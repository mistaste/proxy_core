package middleware

import (
	"segment/liboutline"
	"segment/libxray"

	"github.com/tidwall/gjson"
)

func DetectCoreNameFromConfig(config string) string {
	
	hasServer := gjson.Get(config, "server").Exists()
	hasServerPort := gjson.Get(config, "server_port").Exists()
	hasMethod := gjson.Get(config, "method").Exists()
	hasPassword := gjson.Get(config, "password").Exists()

	if hasServer && hasServerPort && hasMethod && hasPassword {
		return liboutline.GetOutlineService().CoreName()
	}

	return libxray.GetXrayService().CoreName()
}
