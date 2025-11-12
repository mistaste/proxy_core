package ios

import (
	"context"
	"fmt"
	"strings"

	"segment/proxycoreproto"
	"segment/server"
)


func StartGRPCIOS() bool {
	return server.StartGRPCServer()
}



func StartCoreIOS(coreName string, dir string, config string, memory int32, isString bool, proxyPort int32) string {
	ctx := context.Background()

	req := &proxycoreproto.StartCoreRequest{
		CoreName:  coreName,
		Dir:       dir,
		Config:    config,
		Memory:    memory,
		IsString:  isString,
		ProxyPort: proxyPort,
		IsVpnMode: false,
	}

	_, err := server.HandleStartCore(ctx, req)
	if err != nil {
		return "ERROR_CORE: " + err.Error()
	}
	return "true"
}


func StopCoreIOS() bool {
	ctx := context.Background()
	_, err := server.HandleStopCore(ctx, &proxycoreproto.Empty{})
	return err == nil
}


func IsCoreRunningIOS() bool {
	ctx := context.Background()
	resp, err := server.HandleIsCoreRunning(ctx, &proxycoreproto.Empty{})
	if err != nil {
		return false
	}
	return resp.Message
}


func GetVersionIOS() string {
	ctx := context.Background()
	resp, err := server.HandleGetVersion(ctx, &proxycoreproto.Empty{})
	if err != nil {
		return "unknown"
	}
	return resp.Message
}


func MeasurePingIOS(urls string) string {
	ctx := context.Background()

	urlList := []string{}
	for _, u := range strings.Split(urls, ",") {
		u = strings.TrimSpace(u)
		if u != "" {
			urlList = append(urlList, u)
		}
	}

	req := &proxycoreproto.MeasurePingRequest{Url: urlList}
	resp, err := server.HandleMeasurePing(ctx, req)
	if err != nil {
		return err.Error()
	}

	delays := []string{}
	for _, r := range resp.Results {
		delays = append(delays, fmt.Sprintf("%d", r.Delay))
	}
	return strings.Join(delays, ",")
}


func FetchLogsIOS() string {
	ctx := context.Background()
	resp, err := server.HandleFetchLogs(ctx, &proxycoreproto.Empty{})
	if err != nil {
		return ""
	}
	return resp.Logs
}


func ClearLogsIOS() {
	ctx := context.Background()
	_, _ = server.HandleClearLogs(ctx, &proxycoreproto.Empty{})
}
