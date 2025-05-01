package util

import (
    "fmt"
    "github.com/fatih/color"
    "sync"
)

var (
    colorsCache   = make(map[color.Attribute]*color.Color)
    mu sync.Mutex
)

func Colorize(c color.Attribute, v ...any) string {
    if !inBetween(c, color.FgBlack, color.FgWhite) {
        return fmt.Sprint(v...)
    }

    return getCachedColor(c).Sprint(v...)
}

func Colorizef(c color.Attribute, format string, v ...any) string {
    if !inBetween(c, color.FgBlack, color.FgWhite) {
        return fmt.Sprintf(format, v...)
    }

    return getCachedColor(c).Sprintf(format, v...)
}

func Colorizeln(c color.Attribute, v ...any) string {
    if !inBetween(c, color.FgBlack, color.FgWhite) {
        return fmt.Sprintln(v...)
    }

    return getCachedColor(c).Sprintln(v...)
}

func getCachedColor(p color.Attribute) *color.Color {
    mu.Lock()
    defer mu.Unlock()

    c, ok := colorsCache[p]
    if !ok {
        c = color.New(p)
        colorsCache[p] = c
    }

    return c
}

func inBetween(i, min, max color.Attribute) bool {
    return (i >= min) && (i <= max)
}
