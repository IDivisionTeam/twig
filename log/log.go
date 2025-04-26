package log

import (
    "github.com/fatih/color"
    "log"
    "twig/util"
)

type Recorder interface {
    Print(v ...any)
    Printf(format string, v ...any)
    Println(v ...any)
}

type TwigRecorder struct {
    level  Level
    logger *log.Logger
    color  color.Attribute
}

type TwigExceptionRecorder struct {
    level  Level
    logger *log.Logger
}

type noOpTwigRecorder struct{}

func (tr *TwigRecorder) Print(v ...any) {
    tr.logger.Print(util.Colorize(tr.color, v...))
}

func (tr *TwigRecorder) Printf(format string, v ...any) {
    tr.logger.Print(util.Colorizef(tr.color, format, v...))
}

func (tr *TwigRecorder) Println(v ...any) {
    tr.logger.Print(util.Colorizeln(tr.color, v...))
}

func (ter *TwigExceptionRecorder) Print(v ...any) {
    if ter.level == FatalLevel {
        ter.logger.Fatal(util.Colorize(color.FgRed, v...))
    } else if ter.level == PanicLevel {
        ter.logger.Panic(util.Colorize(color.FgRed, v...))
    }
}

func (ter *TwigExceptionRecorder) Printf(format string, v ...any) {
    if ter.level == FatalLevel {
        ter.logger.Fatal(util.Colorizef(color.FgRed, format, v...))
    } else if ter.level == PanicLevel {
        ter.logger.Panic(util.Colorizef(color.FgRed, format, v...))
    }
}

func (ter *TwigExceptionRecorder) Println(v ...any) {
    if ter.level == FatalLevel {
        ter.logger.Fatal(util.Colorizeln(color.FgRed, v...))
    } else if ter.level == PanicLevel {
        ter.logger.Panic(util.Colorizeln(color.FgRed, v...))
    }
}

func (i *noOpTwigRecorder) Print(v ...any) {
    // no-op
}

func (i *noOpTwigRecorder) Printf(format string, v ...any) {
    // no-op
}

func (i *noOpTwigRecorder) Println(v ...any) {
    // no-op
}
