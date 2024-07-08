package recorder

import (
    "log"
    "os"
    "sync"
)

type Type int

const (
    INFO Type = iota
    WARN
    ERROR
)

var (
    loggers = make(map[Type]Recorder)
    once    sync.Once
)

func createRecorders() {
    loggers[INFO] = newInfoRecorder()
    loggers[WARN] = newWarningRecorder()
    loggers[ERROR] = newErrorRecorder()
}

func newInfoRecorder() Recorder {
    return &InfoRecorder{
        logger: log.New(os.Stdout, "INFO: ", log.Lmsgprefix|log.LstdFlags),
    }
}

func newWarningRecorder() Recorder {
    return &WarningRecorder{
        logger: log.New(os.Stdout, "WARNING: ", log.Lmsgprefix|log.LstdFlags),
    }
}

func newErrorRecorder() Recorder {
    return &ErrorRecorder{
        logger: log.New(os.Stderr, "ERROR: ", log.Lmsgprefix|log.LstdFlags|log.Lshortfile),
    }
}

func Print(lt Type, v ...any) {
    once.Do(createRecorders)
    loggers[lt].print(v...)
}

func Printf(lt Type, format string, v ...any) {
    once.Do(createRecorders)
    loggers[lt].printf(format, v...)
}

func Println(lt Type, v ...any) {
    once.Do(createRecorders)
    loggers[lt].println(v...)
}
