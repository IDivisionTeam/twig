package log

import (
    "github.com/fatih/color"
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
    PanicLevel
    FatalLevel
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
    loggers[PanicLevel] = newPanicRecorder()
    loggers[FatalLevel] = newFatalRecorder()
}

func newDebugRecorder() Recorder {
    if currentMinLevel <= DebugLevel {
        return &TwigRecorder{
            level:  DebugLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgGreen,
        }
    }
    return &noOpTwigRecorder{}
}

func newInfoRecorder() Recorder {
    if currentMinLevel <= InfoLevel {
        return &TwigRecorder{
            level:  InfoLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgBlue,
        }
    }
    return &noOpTwigRecorder{}
}

func newWarningRecorder() Recorder {
    if currentMinLevel <= WarnLevel {
        return &TwigRecorder{
            level:  WarnLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgYellow,
        }
    }
    return &noOpTwigRecorder{}
}

func newErrorRecorder() Recorder {
    if currentMinLevel <= ErrorLevel {
        return &TwigRecorder{
            level:  ErrorLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
            color:  color.FgRed,
        }
    }
    return &noOpTwigRecorder{}
}

func newPanicRecorder() Recorder {
    if currentMinLevel <= PanicLevel {
        return &TwigExceptionRecorder{
            level:  PanicLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
        }
    }
    return &noOpTwigRecorder{}
}

func newFatalRecorder() Recorder {
    if currentMinLevel <= FatalLevel {
        return &TwigExceptionRecorder{
            level:  FatalLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
        }
    }
    return &noOpTwigRecorder{}
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

func Panic() Recorder {
    once.Do(createRecorders)
    return loggers[PanicLevel]
}

func Fatal() Recorder {
    once.Do(createRecorders)
    return loggers[FatalLevel]
}
