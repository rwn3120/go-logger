package logger

import (
    "testing"
    "fmt"
)

var levels = []string{"ERROR", "WARN", "INFO", "DEBUG", "TRACE"}

func TestGetLevel(t *testing.T) {
    if GetLevel("Error") != LogError {
        t.Error("wrong level")
    }
    if GetLevel("wArn") != LogWarn {
        t.Error("wrong level")
    }
    if GetLevel("infO") != LogInfo {
        t.Error("wrong level")
    }
    if GetLevel("DEBUG") != LogDebug {
        t.Error("wrong level")
    }
    if GetLevel("trace") != LogTrace {
        t.Error("wrong level")
    }
    if GetLevel("blabla") != LogError{
        t.Error("wrong level")
    }
}

func TestLogger(t *testing.T) {
    for i, level := range levels {
        name := fmt.Sprintf("logger-%d", i)
        logger, err := New(name, &Configuration{Level: level})
        if err != nil {
            t.Error("unexpected error:", err.Error())
        }
        if logger.name != name {
            t.Error("unexpected name:", logger.name)
        }
        if logger.level != GetLevel(level) {
            t.Error("unexpected level:", logger.level)
        }
    }
}
func TestLevel(t *testing.T) {
    if _, err := New("logger", &Configuration{Level: ""}); err == nil {
        t.Error("error is nil")
    }

    if _, err := New("logger", &Configuration{Level: "AA"}); err == nil {
        t.Error("error is nil")
    }

    for _, level := range levels {
        if _, err := New("logger", &Configuration{Level: level}); err != nil {
            t.Error(level,"- unexpected error: ", err.Error())
        }
    }
}
