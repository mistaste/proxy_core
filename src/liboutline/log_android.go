

package liboutline








import "C"
import (
	"io"
	"log/slog"
	"os"
	"unsafe"
)

type androidLogWriter struct{}

func (androidLogWriter) Write(p []byte) (n int, err error) {
	cstr := C.CString(string(p))
	C.__android_log_write(C.ANDROID_LOG_INFO, C.CString("outline"), cstr)
	C.free(unsafe.Pointer(cstr))
	return len(p), nil
}


func (osrv *OutlineService) initLogger() {
	level := slog.LevelDebug

	w := io.MultiWriter(os.Stdout, osrv.logWriter, androidLogWriter{})
	h := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	})
	osrv.logger = slog.New(h)
}
