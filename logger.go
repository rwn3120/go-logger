package logger

import (
    "fmt"
    "os"
    "log"
    "strings"
    "runtime"
    "github.com/rwn3120/go-conf"
    "github.com/rwn3120/go-multierror"
    "errors"
)

type Level int

const (
    LogFatal Level = 1 << iota
    LogError
    LogWarn
    LogInfo
    LogDebug
    LogTrace

    LogFatalText = "FATAL"
    LogErrorText = "ERROR"
    LogWarnText  = "WARN"
    LogInfoText  = "INFO"
    LogDebugText = "DEBUG"
    LogTraceText = "TRACE"
)

type Configuration struct {
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

func (c *Configuration) Validate() []error {
    me := multierror.New()
    switch strings.ToUpper(c.Level) {
    case LogErrorText:
    case LogWarnText:
    case LogInfoText:
    case LogDebugText:
    case LogTraceText:
    default:
        me.Add(errors.New(fmt.Sprintf("unknown log level: %s", c.Level)))
    }
    return me.ErrorsOrNil()
}

func (c *Configuration) StdOut() *os.File {
    if c.stdout != nil {
        return c.stdout
    }
    c.stdout = openLog(c.StdOutFile, os.Stdout)
    return c.stdout
}

func (c *Configuration) StdErr() *os.File {
    if c.stderr != nil {
        return c.stderr
    }
    c.stderr = openLog(c.StdErrFile, os.Stderr)
    return c.stderr
}

type Logger struct {
    name  string
    level Level
    fatLogger,
    errLogger,
    wrnLogger,
    infLogger,
    dbgLogger,
    trcLogger *log.Logger
    configuration *Configuration
}

func createPrefix(level string, name string) string {
    return fmt.Sprintf("%6s [%s]", level, name)
}

func GetLevel(level string) Level {
    logLevel := LogError
    switch strings.ToUpper(level) {
    case LogErrorText:
        logLevel = LogError
    case LogWarnText:
        logLevel = LogWarn
    case LogInfoText:
        logLevel = LogInfo
    case LogDebugText:
        logLevel = LogDebug
    case LogTraceText:
        logLevel = LogTrace
    default:
        fmt.Errorf("unknown log level: %s (will use %s instead)", level, LogErrorText)
        logLevel = LogError
    }
    return logLevel
}

func New(name string, configuration *Configuration, flags ...int) (*Logger, error) {
    if !conf.IsValid(configuration) {
        return nil, errors.New("configuration is not valid")
    }

    if len(strings.TrimSpace(name)) == 0 {
        name = "UnnamedLogger"
    }

    flag := 0
    if len(flags) == 0 {
        flag = log.Ldate | log.Ltime
    } else {
        for f := range flags {
            flag = flag | f
        }
    }

    return &Logger{
        name:          name,
        level:         GetLevel(configuration.Level),
        fatLogger:     log.New(configuration.StdErr(), createPrefix(LogFatalText, name), flag),
        errLogger:     log.New(configuration.StdErr(), createPrefix(LogErrorText, name), flag),
        wrnLogger:     log.New(configuration.StdErr(), createPrefix(LogWarnText, name), flag),
        infLogger:     log.New(configuration.StdOut(), createPrefix(LogInfoText, name), flag),
        dbgLogger:     log.New(configuration.StdOut(), createPrefix(LogDebugText, name), flag),
        trcLogger:     log.New(configuration.StdOut(), createPrefix(LogTraceText, name), flag),
        configuration: configuration,
    }, nil
}

func (jl *Logger) log(requiredLevel Level, logger *log.Logger, format string, args ...interface{}) {
    if requiredLevel <= jl.level {
        if jl.level >= LogTrace {
            _, fn, line, _ := runtime.Caller(2)
            prefix := fmt.Sprintf("%s:%d", fn, line)
            logger.Println(fmt.Sprintf(prefix+" "+format, args...))
        } else {
            logger.Println(fmt.Sprintf(format, args...))
        }

    }
}

func (jl *Logger) Fatal(format string, args ...interface{}) {
    jl.log(LogFatal, jl.errLogger, format, args...)
    panic(fmt.Sprintf(format, args...))
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
