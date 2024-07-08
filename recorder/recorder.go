package recorder

import (
    "log"
)

type Recorder interface {
    print(v ...any)
    printf(format string, v ...any)
    println(v ...any)
}

type InfoRecorder struct {
    logger *log.Logger
}

type WarningRecorder struct {
    logger *log.Logger
}

type ErrorRecorder struct {
    logger *log.Logger
}

func (i *InfoRecorder) print(v ...any) {
    i.logger.Print(v...)
}

func (i *InfoRecorder) printf(format string, v ...any) {
    i.logger.Printf(format, v...)
}

func (i *InfoRecorder) println(v ...any) {
    i.logger.Println(v...)
}

func (i *WarningRecorder) print(v ...any) {
    i.logger.Print(v...)
}

func (i *WarningRecorder) printf(format string, v ...any) {
    i.logger.Printf(format, v...)
}

func (i *WarningRecorder) println(v ...any) {
    i.logger.Println(v...)
}

func (i *ErrorRecorder) print(v ...any) {
    i.logger.Print(v...)
}

func (i *ErrorRecorder) printf(format string, v ...any) {
    i.logger.Printf(format, v...)
}

func (i *ErrorRecorder) println(v ...any) {
    i.logger.Println(v...)
}
