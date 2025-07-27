package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-stack/stack"
	"github.com/sirupsen/logrus"
)

type severity string

const (
	severityDebug    severity = "DEBUG"
	severityInfo     severity = "INFO"
	severityWarning  severity = "WARNING"
	severityError    severity = "ERROR"
	severityCritical severity = "CRITICAL"
	severityAlert    severity = "ALERT"
)

var levelsToSeverity = map[logrus.Level]severity{
	logrus.DebugLevel: severityDebug,
	logrus.InfoLevel:  severityInfo,
	logrus.WarnLevel:  severityWarning,
	logrus.ErrorLevel: severityError,
	logrus.FatalLevel: severityCritical,
	logrus.PanicLevel: severityAlert,
}

type sourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

type entry map[string]interface{}

// Formatter implements Stackdriver formatting for logrus.
type Formatter struct {
	Service                string
	Version                string
	StackSkip              []string
	IncludeTransactionName bool
}

// Option lets you configure the Formatter.
type Option func(*Formatter)

// CustomFormatter is a logrus formatter that includes transaction name
type CustomFormatter struct {
	IncludeTransactionName bool // Flag to include transaction name
}

// Format formats the log entry for CustomFormatter
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Prepare the log message
	logMessage := entry.Message

	// Include transaction name if the flag is set
	if f.IncludeTransactionName {
		if transactionName, ok := entry.Data["transaction_name"]; ok {
			logMessage = logMessage + " | Transaction: " + transactionName.(string)
		}
	}

	// Format the log entry as needed (e.g., JSON, plain text)
	return []byte(logMessage + "\n"), nil
}

// WithService lets you configure the service name used for error reporting.
func WithService(n string) Option {
	return func(f *Formatter) {
		f.Service = n
	}
}

// WithVersion lets you configure the service version used for error reporting.
func WithVersion(v string) Option {
	return func(f *Formatter) {
		f.Version = v
	}
}

// WithStackSkip lets you configure which packages should be skipped for locating the error.
func WithStackSkip(v string) Option {
	return func(f *Formatter) {
		f.StackSkip = append(f.StackSkip, v)
	}
}

// NewFormatter returns a new Formatter.
func NewFormatter(options ...Option) *Formatter {
	fmtr := Formatter{
		StackSkip: []string{
			"github.com/sirupsen/logrus",
		},
	}
	for _, option := range options {
		option(&fmtr)
	}
	return &fmtr
}

func (f *Formatter) errorOrigin() (stack.Call, error) {
	skip := func(pkg string) bool {
		for _, skip := range f.StackSkip {
			if pkg == skip {
				return true
			}
		}
		return false
	}

	// We start at 2 to skip this call and our caller's call.
	for i := 2; ; i++ {
		c := stack.Caller(i)
		// ErrNoFunc indicates we're over traversing the stack.
		if _, err := c.MarshalText(); err != nil {
			return stack.Call{}, nil
		}
		pkg := fmt.Sprintf("%+k", c)
		// Remove vendoring from package path.
		parts := strings.SplitN(pkg, "/vendor/", 2)
		pkg = parts[len(parts)-1]
		if !skip(pkg) {
			return c, nil
		}
	}
}

// WithTransactionName is a custom option to include transaction name in the log entry
func WithTransactionName() Option {
	return func(f *Formatter) {
		f.IncludeTransactionName = true // Flag to include transaction name
	}
}

// getColorByLevel returns the color code for a given log level
func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorBlue
	case logrus.InfoLevel:
		return colorGreen
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorGray
	}
}

// Format formats a logrus entry according to the Stackdriver specifications.
func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	severity := levelsToSeverity[e.Level]

	ee := entry{
		LoggerField().Message:  e.Message,
		LoggerField().Severity: severity,
	}

	ee[LoggerField().Timestamp] = time.Now().UTC().Format(time.RFC3339)

	// Extract all custom field and push into log
	for k, v := range e.Data {
		ee[k] = v
	}

	switch severity {
	case severityError, severityCritical, severityAlert:
		if err, ok := e.Data["error"]; ok {
			ee["error"] = fmt.Sprintf("%s: %s", e.Message, err)
		} else {
			ee["error"] = e.Message
		}

		if c, err := f.errorOrigin(); err == nil {
			lineNumber, _ := strconv.ParseInt(fmt.Sprintf("%d", c), 10, 64)

			ee[LoggerField().SourceLocation] = &sourceLocation{
				File:     fmt.Sprintf("%+s", c),
				Line:     int(lineNumber),
				Function: fmt.Sprintf("%n", c),
			}
		}
	}

	// Add color to the output if not in production
	if os.Getenv("ENV") != "prod" {
		color := getColorByLevel(e.Level)
		jsonFormatter := &logrus.JSONFormatter{}
		b, err := jsonFormatter.Format(e)
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m\n", color, string(b))), nil
	}

	jsonFormatter := &logrus.JSONFormatter{}
	b, err := jsonFormatter.Format(e)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// ColorFormatter wraps an existing formatter and adds color
type ColorFormatter struct {
	WrappedFormatter logrus.Formatter
}

// Format implements the logrus.Formatter interface
func (f *ColorFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// First, get the output from the wrapped formatter
	formatted, err := f.WrappedFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	// Add color based on log level
	color := getColorByLevel(entry.Level)

	// Create the colored output
	coloredOutput := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, string(formatted))

	return []byte(coloredOutput), nil
}

// NewColorFormatter creates a new color formatter that wraps an existing formatter
func NewColorFormatter(wrappedFormatter logrus.Formatter) *ColorFormatter {
	return &ColorFormatter{
		WrappedFormatter: wrappedFormatter,
	}
}
