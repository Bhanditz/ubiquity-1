package logutil

import (
    "fmt"
    "github.com/op/go-logging"
)

const (
    traceEnter = "ENTER"
    traceExit = "EXIT"
)

func (logger *impLogger) Debug(str string, args ...Args) {
    logger.goLogger.Debugf(str + " %v", args)
}

func (logger *impLogger) Info(str string, args ...Args) {
    logger.goLogger.Infof(str + " %v", args)
}

func (logger *impLogger) Error(str string, args ...Args) {
    logger.goLogger.Errorf(str + " %v", args)
}

func (logger *impLogger) Trace(level logging.Level) func() {
    switch level {
    case DEBUG:
        logger.goLogger.Debug(traceEnter)
        return func() { logger.goLogger.Debug(traceExit) }
    case INFO:
        logger.goLogger.Info(traceEnter)
        return func() { logger.goLogger.Info(traceExit) }
    case ERROR:
        logger.goLogger.Error(traceEnter)
        return func() { logger.goLogger.Error(traceExit) }
    default:
        panic("unknown level")
    }
}


func (param Param) String() string {
    return fmt.Sprintf("{" + param.Name + "=%v}", param.Value)
}
