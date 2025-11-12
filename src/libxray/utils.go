package libxray

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)


func SetEnv(ctx context.Context, dir string) error {
	select {
	case <-ctx.Done():
		return ctx.Err() 
	default:
	}

	err := os.Setenv("xray.location.asset", dir)
	if err != nil {
		return err
	}
	return nil
}



func FreeMemory(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return 
			case <-ticker.C:
				debug.FreeOSMemory()
			}
		}
	}()
}




func MaxMemory(ctx context.Context, value int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err() 
	default:
	}

	
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(value * 1024 * 1024) 

	
	FreeMemory(ctx, 1*time.Second)

	return nil
}


func FreeOSMemory(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err() 
	default:
	}

	debug.FreeOSMemory()
	return nil
}



func replaceInboundSocksPort(config string, port int32) (string, error) {
	
	
	if !gjson.Valid(config) {
		return "", fmt.Errorf("failed: invalid JSON config provided")
	}

	
	inboundsResult := gjson.Get(config, "inbounds")
	if !inboundsResult.Exists() || !inboundsResult.IsArray() {
		return "", fmt.Errorf("failed: 'inbounds' array not found or incorrect format")
	}

	modifiedConfig := config 

	
	inboundsResult.ForEach(func(key, value gjson.Result) bool {
		
		protocol := value.Get("protocol").String()

		if protocol == "socks" {
			
			
			portPath := fmt.Sprintf("inbounds.%s.port", key.String())

			
			
			
			var err error
			modifiedConfig, err = sjson.Set(modifiedConfig, portPath, port)
			if err != nil {
				
				
				modifiedConfig = "" 
				return false        
			}
			
			return false 
		}
		return true 
	})

	if modifiedConfig == "" {
		
		return "", fmt.Errorf("failed: error during JSON modification to set socks port")
	}

	
	
	
	
	
	
	
	
	
	return modifiedConfig, nil
}
