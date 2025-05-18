package log

import (
    "github.com/fatih/color"
    "log"
    "os"
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

var rp *RecorderPool

type RecorderPool struct {
    minLogLevel Level
    recorders   map[Level]Recorder
}

func init() {
    rp = New()
}

func New() *RecorderPool {
    rp := new(RecorderPool)

    rp.minLogLevel = InfoLevel
    rp.recorders = make(map[Level]Recorder)

    return rp
}

func Reset() {
    rp = New()
}

func SetLevel(l Level) {
    rp.minLogLevel = l
}

func CreateRecorders() {
    rp.createRecorders()
}

func (rp *RecorderPool) createRecorders() {
    rp.recorders[DebugLevel] = newDebugRecorder()
    rp.recorders[InfoLevel] = newInfoRecorder()
    rp.recorders[WarnLevel] = newWarningRecorder()
    rp.recorders[ErrorLevel] = newErrorRecorder()
    rp.recorders[PanicLevel] = newPanicRecorder()
    rp.recorders[FatalLevel] = newFatalRecorder()
}

func newDebugRecorder() Recorder {
    if rp.minLogLevel <= DebugLevel {
        return &TwigRecorder{
            level:  DebugLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgGreen,
        }
    }
    return new(noOpTwigRecorder)
}

func newInfoRecorder() Recorder {
    if rp.minLogLevel <= InfoLevel {
        return &TwigRecorder{
            level:  InfoLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgBlue,
        }
    }
    return new(noOpTwigRecorder)
}

func newWarningRecorder() Recorder {
    if rp.minLogLevel <= WarnLevel {
        return &TwigRecorder{
            level:  WarnLevel,
            logger: log.New(os.Stdout, "", log.Lmsgprefix),
            color:  color.FgYellow,
        }
    }
    return new(noOpTwigRecorder)
}

func newErrorRecorder() Recorder {
    if rp.minLogLevel <= ErrorLevel {
        return &TwigRecorder{
            level:  ErrorLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
            color:  color.FgRed,
        }
    }
    return new(noOpTwigRecorder)
}

func newPanicRecorder() Recorder {
    if rp.minLogLevel <= PanicLevel {
        return &TwigExceptionRecorder{
            level:  PanicLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
        }
    }
    return new(noOpTwigRecorder)
}

func newFatalRecorder() Recorder {
    if rp.minLogLevel <= FatalLevel {
        return &TwigExceptionRecorder{
            level:  FatalLevel,
            logger: log.New(os.Stderr, "", log.Lmsgprefix),
        }
    }
    return new(noOpTwigRecorder)
}

func Print(l Level, v ...any) {
    rp.recorders[l].Print(v...)
}

func Printf(l Level, format string, v ...any) {
    rp.recorders[l].Printf(format, v...)
}

func Println(l Level, v ...any) {
    rp.recorders[l].Println(v...)
}

func Debug() Recorder {
    return rp.recorders[DebugLevel]
}

func Info() Recorder {
    return rp.recorders[InfoLevel]
}

func Warn() Recorder {
    return rp.recorders[WarnLevel]
}

func Error() Recorder {
    return rp.recorders[ErrorLevel]
}

func Panic() Recorder {
    return rp.recorders[PanicLevel]
}

func Fatal() Recorder {
    return rp.recorders[FatalLevel]
}
