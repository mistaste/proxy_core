package log

import (
	"go.uber.org/zap"
)


var Must = zap.Must


type (
	Logger        = zap.Logger
	SugaredLogger = zap.SugaredLogger
)

type (
	
	Option = zap.Option
)


var pkgCallerSkip = zap.AddCallerSkip(2)
