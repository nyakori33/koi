package log

import (
	"bytes"
	"fmt"
	"time"
)

// 日志等级
var level string

type color struct {
	Info  string
	Debug string
	Warn  string
	Error string
	Fatal string
}

// 颜色
var Color = color{
	Info:  "\x1b[37m",
	Debug: "\x1b[32m",
	Warn:  "\x1b[33m",
	Error: "\x1b[91m",
	Fatal: "\x1b[91m",
}

func (*color) Reset() string { return "\x1b[0m" }

func (color *color) Level() string {
	switch level {
	case "INFO":
		return color.Info
	case "DEBUG":
		return color.Debug
	case "WARN":
		return color.Warn
	case "ERROR":
		return color.Error
	case "FATAL":
		return color.Fatal
	default:
		return color.Info
	}
}

// 日志格式
var Format = func(a ...any) string {
	var buffer bytes.Buffer

	buffer.WriteString(Color.Level())
	buffer.WriteString("[")
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	buffer.WriteString("]")
	buffer.WriteString(" ")
	buffer.WriteString("[")
	buffer.WriteString(level)
	buffer.WriteString("]")
	buffer.WriteString(":")
	buffer.WriteString(" ")
	for _, v := range a {
		buffer.WriteString(" ")
		buffer.WriteString(fmt.Sprint(v))
	}
	buffer.WriteString(Color.Reset())

	return buffer.String()
}

func Info(a ...any) {
	level = "INFO"
	fmt.Println(Format(a...))
}

func Debug(a ...any) {
	level = "DEBUG"
	fmt.Println(Format(a...))
}

func Warn(a ...any) {
	level = "WARN"
	fmt.Println(Format(a...))
}

func Error(a ...any) {
	level = "ERROR"
	fmt.Println(Format(a...))
}

func Fatal(a ...any) {
	level = "FATAL"
	fmt.Println(Format(a...))
	panic(fmt.Sprint(a...))
}
