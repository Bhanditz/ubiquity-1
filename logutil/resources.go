package logutil

import (
    "github.com/op/go-logging"
)

const (
    DEBUG = logging.DEBUG
    INFO = logging.INFO
    ERROR = logging.ERROR
)


type Param struct {
    Name string
    Value interface{}
}

type Args []Param

type Logger interface {
    Debug(str string, args ...Args)
    Info(str string, args ...Args)
    Error(str string, args ...Args)
    Trace(level logging.Level) func()
}

type impLogger struct {
    goLogger    *logging.Logger
}