package log

import (
    "log"
    "os"
    "sync"
)

type Level int

const (
    DebugLevel Level = iota - 1
    InfoLevel        // InfoLevel is the default logging priority.
    WarnLevel
    ErrorLevel
)

var currentMinLevel = InfoLevel

func SetLevel(l Level) {
    currentMinLevel = l
}

var (
    loggers = make(map[Level]Recorder)
    once    sync.Once
)

func createRecorders() {
    loggers[DebugLevel] = newDebugRecorder()
    loggers[InfoLevel] = newInfoRecorder()
    loggers[WarnLevel] = newWarningRecorder()
    loggers[ErrorLevel] = newErrorRecorder()
}

func newDebugRecorder() Recorder {
    if currentMinLevel >= DebugLevel {
        return &DebugRecorder{
            logger: log.New(os.Stdout, "DEBUG: ", log.Lmsgprefix),
        }
    }
    return &noOpRecorder{}
}

func newInfoRecorder() Recorder {
    if currentMinLevel >= InfoLevel {
        return &InfoRecorder{
            logger: log.New(os.Stdout, "INFO: ", log.Lmsgprefix),
        }
    }
    return &noOpRecorder{}
}

func newWarningRecorder() Recorder {
    if currentMinLevel >= WarnLevel {
        return &WarningRecorder{
            logger: log.New(os.Stdout, "WARNING: ", log.Lmsgprefix),
        }
    }
    return &noOpRecorder{}
}

func newErrorRecorder() Recorder {
    if currentMinLevel >= ErrorLevel {
        return &ErrorRecorder{
            logger: log.New(os.Stderr, "ERROR: ", log.Lmsgprefix),
        }
    }
    return &noOpRecorder{}
}

func Print(l Level, v ...any) {
    once.Do(createRecorders)
    loggers[l].Print(v...)
}

func Printf(l Level, format string, v ...any) {
    once.Do(createRecorders)
    loggers[l].Printf(format, v...)
}

func Println(l Level, v ...any) {
    once.Do(createRecorders)
    loggers[l].Println(v...)
}

func Debug() Recorder {
    once.Do(createRecorders)
    return loggers[DebugLevel]
}

func Info() Recorder {
    once.Do(createRecorders)
    return loggers[InfoLevel]
}

func Warn() Recorder {
    once.Do(createRecorders)
    return loggers[WarnLevel]
}

func Error() Recorder {
    once.Do(createRecorders)
    return loggers[ErrorLevel]
}
