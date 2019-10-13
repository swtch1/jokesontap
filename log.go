package jokesontap

import (
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

var logTimestampFmt = time.RFC3339

// InitLogger sets values on the global logger for use everywhere else in the application.
func InitLogger(w io.Writer, level string, format string, prettyJson bool) {
	var lvl log.Level
	switch strings.ToLower(level) {
	case "trace":
		lvl = log.TraceLevel
	case "debug":
		lvl = log.DebugLevel
	case "info":
		lvl = log.InfoLevel
	case "warn":
		lvl = log.WarnLevel
	case "error":
		lvl = log.ErrorLevel
	case "fatal":
		lvl = log.FatalLevel
	default:
		log.Fatalf("unexpected level '%s'", level)
	}

	var formatter log.Formatter
	switch strings.ToLower(format) {
	case "json":
		formatter = &log.JSONFormatter{TimestampFormat: logTimestampFmt, PrettyPrint: prettyJson}
	case "text":
		formatter = &log.TextFormatter{TimestampFormat: logTimestampFmt}
	default:
		formatter = &log.TextFormatter{TimestampFormat: logTimestampFmt}
		log.Fatalf("unexpected format '%s'", format)

	}
	log.SetLevel(lvl)
	log.SetFormatter(formatter)
	log.SetOutput(w)
	log.SetReportCaller(true)
}
