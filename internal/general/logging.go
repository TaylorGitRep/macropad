package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func padRightSide(str string, item string, count int) string {
	return str + strings.Repeat(item, count)
}

func getTimezoneOffset() string {
	t := time.Now()
	_, offset := t.Zone()
	var realOffset string
	if offset == 0 {
		realOffset = "Z"
	} else {
		realOffset = fmt.Sprintf("%03d", offset/60/60)
		realOffset = fmt.Sprintf("%v", padRightSide(realOffset, "0", 2))
	}
	return realOffset
}

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg%"
	defaultTimestampFormat = time.RFC3339
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Timestamp format
	TimestampFormat string
	// Available standard keys: time, msg, lvl
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat         string
	LogLevelPadding   bool
	TimestampTimezone bool
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if f.TimestampTimezone {
		output = strings.Replace(output, "%time%", "%time%"+getTimezoneOffset(), 1)
	}
	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	// Remove new lines from message and trim space.
	msg := entry.Message
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\n", "")
	msg = strings.TrimSpace(msg)

	output = strings.Replace(output, "%msg%", msg, 1)

	level := strings.ToUpper(entry.Level.String())
	if f.LogLevelPadding {
		level = fmt.Sprintf("%-7s", level)
	}
	output = strings.Replace(output, "%lvl%", level, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}

func LoggingInit() {
	logFolder := fmt.Sprintf("%v/log/", GetAppStorageFolder())
	var _ = os.MkdirAll(logFolder, os.ModePerm)
	filename := filepath.Base(os.Args[0])
	filename = strings.ReplaceAll(filename, filepath.Ext(filename), "")
	logloc := fmt.Sprintf("%v/%v.log", logFolder, filename)

	f, err := os.OpenFile(logloc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to initialize log file due to the following error: %v", err.Error())
		os.Exit(1)
	}

	f.Close()

	var logLevel logrus.Level
	level := os.Getenv("LOGLEVEL")
	if level == "" {
		level = os.Getenv("loglevel")
	}
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "warning":
		logLevel = logrus.WarnLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "panic":
		logLevel = logrus.PanicLevel
	default:
		logLevel = logrus.InfoLevel
	}

	timezoneFormat := "2006-01-02T15:04:05"

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logloc,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     30, //days
		Level:      logLevel,
		Formatter: &Formatter{
			TimestampFormat:   timezoneFormat,
			TimestampTimezone: true,
			LogFormat:         "(%time%) [%lvl%] %msg%\n",
			LogLevelPadding:   true,
		},
	})
	logrus.AddHook(rotateFileHook)

	if err != nil {
		fmt.Printf("Failed to initialize file rotate hook: %v", err.Error())
		os.Exit(1)
	}

	logrus.SetFormatter(&Formatter{
		LogFormat: "%msg%\n",
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
}
