package log

import (
	"bytes"
	"sync"

	alog "github.com/GFW-knocker/Xray-core/app/log"
	"github.com/GFW-knocker/Xray-core/common"
	"github.com/GFW-knocker/Xray-core/common/log"
)

var (
	logBuffer      = &bytes.Buffer{}
	logMutex       sync.RWMutex
	loggerAdded    bool
	maxBufferSize  = 2 * 1024 * 1024 
	loggingEnabled = true
)


func WriteLogToBuffer(msg string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if !loggingEnabled {
		return
	}

	
	if logBuffer.Len() > maxBufferSize {
		data := logBuffer.Bytes()
		if cut := bytes.IndexByte(data[len(data)/2:], '\n'); cut != -1 {
			logBuffer = bytes.NewBuffer(data[len(data)/2+cut+1:])
		} else {
			logBuffer.Reset()
		}
	}

	logBuffer.WriteString(msg + "\n")
}


func StartLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if loggerAdded {
		return
	}

	logBuffer.Reset()

	common.Must(alog.RegisterHandlerCreator(alog.LogType_Console, func(_ alog.LogType, _ alog.HandlerCreatorOptions) (log.Handler, error) {
		return registerPlatformLogger(), nil
	}))

	loggerAdded = true
}


func StopLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	logBuffer.Reset()
	loggerAdded = false
}


func FetchLogs() string {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logBuffer.Len() == 0 {
		return ""
	}

	logs := logBuffer.String()
	logBuffer.Reset()
	return logs
}

func ClearLogs() bool {
	logMutex.Lock()
	defer logMutex.Unlock()

	loggingEnabled = false
	logBuffer.Reset()
	loggingEnabled = true 
	return true
}
