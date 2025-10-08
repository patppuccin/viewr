package out

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/helpers"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/fatih/color"
)

// Banner /////////////////////////////////

func Banner(message string) string {
	var sb strings.Builder
	sb.WriteString(constants.AppBanner)
	sb.WriteString("\n")
	blue := color.New(color.FgBlue).SprintFunc()
	sb.WriteString(blue(message))
	sb.WriteString("\n")
	return sb.String()
}

// Console Logger ///////////////////////////

type SCLogger struct{}

func (l SCLogger) printLog(icon string, msg string) {
	message := color.New(color.FgWhite).Sprint(msg)
	os.Stdout.WriteString(icon + " " + message + "\n")
}

func (l SCLogger) Debug(msg string) {
	icon := color.New(color.FgCyan, color.Bold).Sprint("[?]")
	l.printLog(icon, msg)
}

func (l SCLogger) Info(msg string) {
	icon := color.New(color.FgBlue, color.Bold).Sprint("[i]")
	l.printLog(icon, msg)
}

func (l SCLogger) Warn(msg string) {
	icon := color.New(color.FgYellow, color.Bold).Sprint("[!]")
	l.printLog(icon, msg)
}

func (l SCLogger) Error(msg string) {
	icon := color.New(color.FgRed, color.Bold).Sprint("[✕]")
	l.printLog(icon, msg)
}

func (l SCLogger) Success(msg string) {
	icon := color.New(color.FgGreen, color.Bold).Sprint("[✓]")
	l.printLog(icon, msg)
}

var Logger = &SCLogger{}

// Structured Logger ///////////////////////

var validLogLevels = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
}

func NewStructuredLogger(logLevel string, logToConsole bool) (*zerolog.Logger, error) {

	// Set global time format for zerolog
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	// Set global log level
	level, ok := validLogLevels[logLevel]
	if !ok {
		level = zerolog.InfoLevel // default
	}
	zerolog.SetGlobalLevel(level)

	rootDir, err := os.Executable()
	if err != nil {
		return nil, helpers.SafeErr("failed to resolve root path", err)
	}
	rootPath := filepath.Dir(rootDir)

	// Create log file writer
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(rootPath, "logs", constants.AppAbbrName+".log.json"),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Setup writer
	var writer io.Writer
	if logToConsole {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFieldFormat,
		}
		writer = zerolog.MultiLevelWriter(consoleWriter, logFile)
	} else {
		writer = logFile
	}

	logger := zerolog.New(writer).Level(level).With().Timestamp().Logger()
	return &logger, nil
}
