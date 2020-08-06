// Logging interface that allows filtering messages by level. Call
// `RegisterDefault()` to register `logger` as the default logger. After
// that, any call to Go's `log` package will get redirected through this
// package.
//
// After registering `logger`, call any of `Debugf()`, Infof()`, Warnf()`,
// Errorf()` or Fatalf()` to specify the level of the message to be logged.

package logger

import (
    "fmt"
    "log"
    "io"
)

// Default log level when `log` is used directly.
type LogLevel uint
const (
    LogDebug LogLevel = iota
    LogInfo
    LogWarn
    LogError
    LogFatal
    LogMax
)

// The internal logger
type _logger struct {
    // Default log level as used by this logger.
    defaultLevel LogLevel
    // Maximum level that may be logged by this logger.
    maxLevel LogLevel
    // Writers used by each log type
    output []*log.Logger
}

// Log a debug message to the regular output.
func (l *_logger) Debugf(fmtStr string, v ...interface{}) {
    if l.maxLevel > LogDebug {
        return
    }
    l.output[LogDebug].Printf(fmtStr, v...)
}

// Log an information message to the regular output.
func (l *_logger) Infof(fmtStr string, v ...interface{}) {
    if l.maxLevel > LogInfo {
        return
    }
    l.output[LogInfo].Printf(fmtStr, v...)
}

// Log a warning message to the error output.
func (l *_logger) Warnf(fmtStr string, v ...interface{}) {
    if l.maxLevel > LogWarn {
        return
    }
    l.output[LogWarn].Printf(fmtStr, v...)
}

// Log an error message to the error output.
func (l *_logger) Errorf(fmtStr string, v ...interface{}) {
    if l.maxLevel > LogError {
        return
    }
    l.output[LogError].Printf(fmtStr, v...)
}

// Log a fatal error message to the error output and `panic`.
func (l *_logger) Fatalf(fmtStr string, v ...interface{}) {
    if l.maxLevel <= LogFatal {
        l.output[LogFatal].Printf(fmtStr, v...)
    }
    panic(fmt.Sprintf(fmtStr, v...))
}

// Implement io.Write to redirect `log` messages.
func (l *_logger) Write(p []byte) (int, error) {
    if l.maxLevel <= l.defaultLevel {
        msg := string(p)
        switch l.defaultLevel {
        case LogDebug:
            l.Debugf(msg)
        case LogInfo:
            l.Infof(msg)
        case LogWarn:
            l.Warnf(msg)
        case LogError:
            l.Errorf(msg)
        case LogFatal:
            l.Fatalf(msg)
        default:
            /* Shouldn't happen... */
            l.Warnf(msg)
        }
    }

    return len(p), nil
}

// Register a `logger` as the default logger.
func RegisterDefault(def LogLevel, max LogLevel, out io.Writer, err io.Writer) {
    if def >= LogMax {
        log.Fatalf("logger: Invalid default log level!")
    } else if max >= LogMax {
        log.Fatalf("logger: Invalid maximum log level!")
    } else if def < max {
        log.Fatalf("logger: Default log level must be smaller than maximum log level!")
    }

    // Create a custom logger for each level
    var output []*log.Logger
    flags := log.Ldate | log.Ltime | log.Lmsgprefix

    for i := LogLevel(0); i < LogMax; i++ {
        switch i {
            case LogDebug:
                l := log.New(out, "[DEBUG] ", flags)
                output = append(output, l)
            case LogInfo:
                l := log.New(out, "[INFO.] ", flags)
                output = append(output, l)
            case LogWarn:
                l := log.New(err, "[WARN.] ", flags)
                output = append(output, l)
            case LogError:
                l := log.New(err, "[ERROR] ", flags)
                output = append(output, l)
            case LogFatal:
                l := log.New(err, "[PANIC] ", flags)
                output = append(output, l)
            default:
                panic(fmt.Sprintf("logger: Unsupported logLevel %+v", LogLevel(i)))
        }
    }

    l := _logger {
        defaultLevel: def,
        maxLevel: max,
        output: output,
    }

    log.SetOutput(&l)
    // Remove the flags, so anything logged by `log.*` won't use the flags
    // twice (once for log's flags and once for logger's flags)
    log.SetFlags(0)
}

// Log a debug message to the regular output. If `logger` hasn't been
// configured as the default `log`, it simply logs to the default.
func Debugf(fmtStr string, v ...interface{}) {
    if l, ok := log.Writer().(*_logger); ok {
        l.Debugf(fmtStr, v...)
    } else {
        log.Printf(fmtStr, v...)
    }
}

// Log an information message to the regular output. If `logger` hasn't
// been configured as the default `log`, it simply logs to the default.
func Infof(fmtStr string, v ...interface{}) {
    if l, ok := log.Writer().(*_logger); ok {
        l.Infof(fmtStr, v...)
    } else {
        log.Printf(fmtStr, v...)
    }
}

// Log a warning message to the error output. If `logger` hasn't been
// configured as the default `log`, it simply logs to the default.
func Warnf(fmtStr string, v ...interface{}) {
    if l, ok := log.Writer().(*_logger); ok {
        l.Warnf(fmtStr, v...)
    } else {
        log.Printf(fmtStr, v...)
    }
}

// Log an error message to the error output. If `logger` hasn't been
// configured as the default `log`, it simply logs to the default.
func Errorf(fmtStr string, v ...interface{}) {
    if l, ok := log.Writer().(*_logger); ok {
        l.Errorf(fmtStr, v...)
    } else {
        log.Printf(fmtStr, v...)
    }
}

// Log a fatal error message to the error output and `panic`. If `logger`
// hasn't been configured as the default `log`, it simply logs to the
// default.
func Fatalf(fmtStr string, v ...interface{}) {
    if l, ok := log.Writer().(*_logger); ok {
        l.Fatalf(fmtStr, v...)
    } else {
        log.Fatalf(fmtStr, v...)
    }
}
