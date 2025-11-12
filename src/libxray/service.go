package libxray

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GFW-knocker/Xray-core/common/cmdarg"
	"github.com/GFW-knocker/Xray-core/core"
	"github.com/GFW-knocker/Xray-core/infra/conf/serial"
)


func (xs *XrayService) loadServer(ctx context.Context, config string, isStr bool, port int32) (*core.Instance, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err() 
	default:
	}

	var jsonConfig *core.Config
	var err error

	if isStr {
		
		config, err = replaceInboundSocksPort(config, port)
		if err != nil {
			return nil, fmt.Errorf("failed: unable to replace inbound socks port: %v", err)
		}

		
		jsonConfig, err = serial.LoadJSONConfig(strings.NewReader(config))
		if err != nil {
			return nil, fmt.Errorf("failed: unable to parse JSON config: %v", err)
		}
	} else {
		
		fileContent, err := os.ReadFile(config)
		if err != nil {
			return nil, fmt.Errorf("failed: unable to read config file: %v", err)
		}

		modifiedConfig, err := replaceInboundSocksPort(string(fileContent), port)
		if err != nil {
			return nil, fmt.Errorf("failed: unable to replace inbound socks port: %v", err)
		}

		
		err = os.WriteFile(config, []byte(modifiedConfig), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed: unable to write modified config file: %v", err)
		}

		
		file := cmdarg.Arg{config}
		jsonConfig, err = core.LoadConfig("json", file)
		if err != nil {
			return nil, fmt.Errorf("failed: unable to load config from file: %v", err)
		}
	}

	
	server, err := core.New(jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("failed: unable to create Xray instance: %v", err)
	}

	return server, nil
}
