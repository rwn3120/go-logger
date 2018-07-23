package gologger

import (
    "fmt"
    "os"
    "log"
    "strings"
    "runtime"
    "github.com/rwn3120/goconf"
)

type LoggingLevel int

const (
    LogError LoggingLevel = 1 << iota
    LogWarn
    LogInfo
    LogDebug
    LogTrace

    LogLevelError = "ERROR"
    LogLevelWarn  = "WARN"
    LogLevelInfo  = "INFO"
    LogLevelDebug = "DEBUG"
    LogLevelTrace = "TRACE"
)

type LoggerConfiguration struct {
    Level      string
    StdOutFile string
    StdErrFile string

    stdout *os.File
    stderr *os.File
}

func openLog(path string, defaultFile *os.File) *os.File {
    if len(strings.TrimSpace(path)) == 0 {
        return defaultFile
    }
    stdOut, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0660)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not open %s: %s\n", path, err.Error())
        return defaultFile
    }
    return stdOut
}

func (c *LoggerConfiguration) Validate() *[]string {
    var errorList []string

    if errorsCount := len(errorList); errorsCount > 0 {
        return &errorList
    }
    return nil
}

func (c *LoggerConfiguration) StdOut() *os.File {
    if c.stdout != nil {
        return c.stdout
    }
    c.stdout = openLog(c.StdOutFile, os.Stdout)
    return c.stdout
}

func (c *LoggerConfiguration) StdErr() *os.File {
    if c.stderr != nil {
        return c.stderr
    }
    c.stderr = openLog(c.StdErrFile, os.Stderr)
    return c.stderr
}

type Logger struct {
    name          string
    level         LoggingLevel
    errLogger     *log.Logger
    wrnLogger     *log.Logger
    infLogger     *log.Logger
    dbgLogger     *log.Logger
    trcLogger     *log.Logger
    configuration *LoggerConfiguration
}

func createPrefix(level string, name string) string {
    return fmt.Sprintf("%6s [%s]", level, name)
}

func NewLogger(name string, configuration *LoggerConfiguration, flags ...int) *Logger {
    if !goconf.IsValid(configuration) {
        panic("LoggerConfiguration is not valid")
    }

    if len(strings.TrimSpace(name)) == 0 {
        name = "UnnamedLogger"
    }

    flag := 0
    if len(flags) == 0 {
        flag = log.Ldate | log.Ltime //| log.Lshortfile
    } else {
        for f := range flags {
            flag = flag | f
        }
    }

    logLevel := LogError
    switch configuration.Level {
    case LogLevelError:
        logLevel = LogError
    case LogLevelWarn:
        logLevel = LogWarn
    case LogLevelInfo:
        logLevel = LogInfo
    case LogLevelDebug:
        logLevel = LogDebug
    case LogLevelTrace:
        logLevel = LogTrace
    default:
        logLevel = LogError
    }
    return &Logger{
        name:          name,
        level:         logLevel,
        errLogger:     log.New(configuration.StdErr(), createPrefix(LogLevelError, name), flag),
        wrnLogger:     log.New(configuration.StdErr(), createPrefix(LogLevelWarn, name), flag),
        infLogger:     log.New(configuration.StdOut(), createPrefix(LogLevelInfo, name), flag),
        dbgLogger:     log.New(configuration.StdOut(), createPrefix(LogLevelDebug, name), flag),
        trcLogger:     log.New(configuration.StdOut(), createPrefix(LogLevelTrace, name), flag),
        configuration: configuration,
    }
}

func (jl *Logger) log(requiredLevel LoggingLevel, logger *log.Logger, format string, args ...interface{}) {
    if jl.level <= requiredLevel {
        if jl.level >= LogDebug {
            _, fn, line, _ := runtime.Caller(2)
            prefix := fmt.Sprintf("%s:%d", fn, line)
            logger.Println(fmt.Sprintf(prefix+" "+format, args...))
        } else {
            logger.Println(fmt.Sprintf(format, args...))
        }

    }
}

func (jl *Logger) Error(format string, args ...interface{}) {
    jl.log(LogError, jl.errLogger, format, args...)
}

func (jl *Logger) Warn(format string, args ...interface{}) {
    jl.log(LogWarn, jl.wrnLogger, format, args...)
}

func (jl *Logger) Debug(format string, args ...interface{}) {
    jl.log(LogDebug, jl.dbgLogger, format, args...)
}

func (jl *Logger) Info(format string, args ...interface{}) {
    jl.log(LogInfo, jl.infLogger, format, args...)
}

func (jl *Logger) Trace(format string, args ...interface{}) {
    jl.log(LogTrace, jl.trcLogger, format, args...)
}
