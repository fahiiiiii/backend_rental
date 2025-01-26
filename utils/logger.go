package utils

import (
    // "fmt"
    "log"
    "os"
    "io"
)

type LoggerConfig struct {
    infoWriter  io.Writer
    warnWriter  io.Writer
    errorWriter io.Writer
}

func NewLogger(opts ...Option) *Logger {
    cfg := &LoggerConfig{
        infoWriter:  os.Stdout,
        warnWriter:  os.Stdout,
        errorWriter: os.Stderr,
    }

    for _, opt := range opts {
        opt(cfg)
    }

    return &Logger{
        infoLog:  log.New(cfg.infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
        warnLog:  log.New(cfg.warnWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
        errorLog: log.New(cfg.errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

type Option func(*LoggerConfig)

func WithInfoOutput(w io.Writer) Option {
    return func(c *LoggerConfig) {
        c.infoWriter = w
    }
}

func WithWarnOutput(w io.Writer) Option {
    return func(c *LoggerConfig) {
        c.warnWriter = w
    }
}

func WithErrorOutput(w io.Writer) Option {
    return func(c *LoggerConfig) {
        c.errorWriter = w
    }
}

type Logger struct {
    infoLog  *log.Logger
    warnLog  *log.Logger
    errorLog *log.Logger
}

func (l *Logger) Info(message string) {
    l.infoLog.Println(message)
}

func (l *Logger) Warn(message string) {
    l.warnLog.Println(message)
}

func (l *Logger) Error(message string, err error) {
    if err != nil {
        l.errorLog.Printf("%s: %v", message, err)
    } else {
        l.errorLog.Println(message)
    }
}