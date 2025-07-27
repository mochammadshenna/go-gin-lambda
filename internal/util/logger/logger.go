package logger

import (
	"ai-service/internal/state"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	Logger  = logrus.New()
	lf      loggerField
	hub     *sentry.Hub
	appName = os.Getenv("ENV")
	tracer  = otel.Tracer("ai-service")
)

const SomeContextKey = contextKey(1)

type contextKey int

// Add this type definition near the top of the file with other types
type ColoredJSONFormatter struct {
	*logrus.JSONFormatter
}

// Add this method for the ColoredJSONFormatter
func (f *ColoredJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Get the standard JSON format
	data, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	// Add color based on log level
	var color int
	switch entry.Level {
	case logrus.DebugLevel:
		color = colorBlue
	case logrus.InfoLevel:
		color = colorGreen
	case logrus.WarnLevel:
		color = colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = colorRed
	default:
		color = colorGray
	}

	// Color the entire JSON output
	return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, string(data))), nil
}

// Init initializes the logger and Sentry with the provided DSN
func Init(sentryDSN string) {
	lf = newLoggerField()

	if appName == "prod" {
		if err := initSentry(sentryDSN); err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		hub = sentry.CurrentHub()
	}
	setLogLevel("debug")
	addSentryHook()
	setFormatter(appName, "1.0")

	if appName != "prod" {
		// Use the ColorFormatter for non-production environments
		Logger.SetFormatter(NewColorFormatter(
			&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
			},
		))
	} else {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}

// initSentry initializes Sentry with enhanced configuration
func initSentry(sentryDSN string) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Transport:        sentry.NewHTTPSyncTransport(),
		Environment:      appName,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampler:    sentry.TracesSampler(func(ctx sentry.SamplingContext) float64 { return 1.0 }),
		// Enhanced configuration for better monitoring
		Debug:            appName == "dev",
		BeforeSend:       beforeSend,
		BeforeBreadcrumb: beforeBreadcrumb,
		Integrations: func(integrations []sentry.Integration) []sentry.Integration {
			// Add custom integrations if needed
			return integrations
		},
	})
}

// beforeSend processes events before sending to Sentry
func beforeSend(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	// Add custom context
	event.Extra["service_name"] = "ai-service"
	event.Extra["version"] = "1.0.0"

	// Filter out certain errors if needed
	if event.Exception != nil {
		for _, exception := range event.Exception {
			if strings.Contains(exception.Value, "connection refused") {
				return nil // Don't send connection errors
			}
		}
	}

	return event
}

// beforeBreadcrumb processes breadcrumbs before sending to Sentry
func beforeBreadcrumb(breadcrumb *sentry.Breadcrumb, hint *sentry.BreadcrumbHint) *sentry.Breadcrumb {
	// Add custom breadcrumb data
	breadcrumb.Data["service"] = "ai-service"
	breadcrumb.Data["timestamp"] = time.Now().Unix()

	return breadcrumb
}

// setLogLevel sets the log level for the logger
func setLogLevel(logConfigLevel string) {
	if logConfigLevel == "" {
		logConfigLevel = "debug"
	}
	logLevel, err := logrus.ParseLevel(logConfigLevel)
	if err != nil {
		log.Fatalf("error setting log level: %v", err)
	}
	Logger.SetLevel(logLevel)
}

// setFormatter sets the custom log formatter
func setFormatter(serviceName, version string) {
	formatter := NewFormatter(
		WithService(serviceName),
		WithVersion(version),
		WithStackSkip("ai-service/internal/util/logger"),
		WithStackSkip("ai-service/internal/util/exception"),
	)

	if appName != "prod" {
		Logger.SetFormatter(&ColoredJSONFormatter{
			JSONFormatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
			},
		})
	} else {
		Logger.SetFormatter(formatter)
	}
}

// addSentryHook adds a Sentry hook to the logger
func addSentryHook() {
	if appName == "prod" {
		hook, err := NewSentryHook()
		if err != nil {
			log.Printf("error adding Sentry hook: %v", err)
			return
		}
		Logger.AddHook(hook)
	}
}

// NewSentryHook creates a new Logrus hook for Sentry
func NewSentryHook() (*SentryHook, error) {
	return &SentryHook{}, nil
}

// SentryHook is a Logrus hook that sends logs to Sentry
type SentryHook struct{}

// Levels returns the log levels that this hook should be triggered for
func (h *SentryHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log event is fired
func (h *SentryHook) Fire(entry *logrus.Entry) error {
	if hub != nil {
		event := &sentry.Event{
			Message: entry.Message,
			Level:   sentry.Level(entry.Level.String()),
			Extra:   entry.Data,
		}

		transactionName := fmt.Sprintf("%s %s", entry.Data["method"], entry.Data["url"])
		event.Transaction = transactionName

		if url, ok := entry.Data["url"]; ok {
			event.Extra["url"] = url
		}
		if traceID, ok := entry.Data["trace_id"]; ok {
			event.Extra["trace_id"] = traceID
		}
		hub.CaptureEvent(event)
	}
	return nil
}

// StartSpanWithFuncName starts a span for tracing with OpenTelemetry
func StartSpanWithFuncName(ctx context.Context, tags map[string]string) {
	segmentName := funcCallerName(3)
	span := sentry.StartSpan(ctx, segmentName)
	span.Sampled = sentry.SampledTrue

	for key, value := range tags {
		span.SetTag(key, value)
	}
	span.Finish()
}

// StartTransactionWithFuncName starts a transaction for the request and automatically finishes it
func StartTransactionWithFuncName(w http.ResponseWriter, r *http.Request) context.Context {
	ctx := r.Context()
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}

	// Get the function name for the transaction
	transactionName := funcCallerName(3)
	moduleName := r.URL.Path

	options := []sentry.SpanOption{
		sentry.WithOpName("ai_app"),
		sentry.ContinueFromRequest(r),
		sentry.WithTransactionSource(sentry.SourceURL),
	}

	// Start the transaction
	transaction := sentry.StartTransaction(ctx, transactionName, options...)

	// Add tags to the transaction
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("function_name", transactionName)
		scope.SetTag("module", moduleName)
	})

	// Log the transaction start
	logrus.WithFields(logrus.Fields{
		"transaction_name": transactionName,
		"url":              r.URL.String(),
		"trace_id":         transaction.TraceID.String(),
	}).Info("Transaction started")

	// Defer the transaction finish to ensure it is called when the function exits
	defer func() {
		transaction.Finish()
	}()

	return transaction.Context()
}

// ErrorHubCaptureFromContext captures errors in Sentry with enhanced context
func ErrorHubCaptureFromContext(ctx context.Context, err error, segmentName string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("segment", segmentName)
			scope.SetExtra("error_message", err.Error())
			scope.SetLevel(sentry.LevelError)

			// Capture the source location of the error
			_, file, line, ok := runtime.Caller(1)
			if ok {
				sourceLocation := fmt.Sprintf("%s:%d", file, line)
				scope.SetExtra("sourceLocation", sourceLocation)
				scope.SetContext("source", map[string]interface{}{
					"file": file,
					"line": line,
				})

				// Log the error with the correct source location
				logrus.WithField("sourceLocation", sourceLocation).Error(err)
			}

			// Add request context if available
			if requestID := ctx.Value(state.HttpHeaders().RequestId); requestID != nil {
				scope.SetTag("request_id", fmt.Sprintf("%v", requestID))
			}

			hub.CaptureException(err)
		})
	}
}

// InfoHubCaptureFromContext captures informational messages in Sentry
func InfoHubCaptureFromContext(c context.Context, message string) {
	if hub := sentry.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelInfo)
			hub.CaptureMessage(message)
		})
	}
}

// loggerField defines the structure of log fields
type loggerField struct {
	RequestId      string `json:"requestId"`
	Latency        string `json:"latency"`
	RequestMethod  string `json:"requestMethod"`
	Resource       string `json:"resource"`
	UserAgent      string `json:"userAgent"`
	PlatformType   string `json:"platformType"`
	Platform       string `json:"platform"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	XForwardedFor  string `json:"xForwardedFor"`
	Message        string `json:"message"`
	Severity       string `json:"severity"`
	Timestamp      string `json:"timestamp"`
	SourceLocation string `json:"sourceLocation"`
}

// newLoggerField initializes a new loggerField
func newLoggerField() loggerField {
	return loggerField{
		RequestId:      "requestId",
		Latency:        "latency",
		RequestMethod:  "requestMethod",
		Resource:       "resource",
		UserAgent:      "userAgent",
		PlatformType:   "platformType",
		Platform:       "platform",
		Version:        "version",
		Status:         "status",
		XForwardedFor:  "xForwardedFor",
		Message:        "message",
		Severity:       "severity",
		Timestamp:      "timestamp",
		SourceLocation: "sourceLocation",
	}
}

// LoggerField returns the current loggerField
func LoggerField() loggerField {
	return lf
}

// logWithFields logs messages with context fields and formatted message
func logWithFields(ctx context.Context, level logrus.Level, format string, args ...interface{}) {
	fields := withFields(ctx)
	message := fmt.Sprintf(format, args...)
	if len(fields) > 0 {
		Logger.WithFields(fields).Log(level, message)
		return
	}
	Logger.Log(level, message)
}

// CaptureSentryEventFromContext captures context information and sends it to Sentry
func CaptureSentryEventFromContext(ctx context.Context, args ...interface{}) {
	if hub == nil {
		return // No Sentry hub available
	}

	event := &sentry.Event{
		Level: sentry.LevelError,
		Extra: logrus.Fields{},
	}

	if requestId := ctx.Value(state.HttpHeaders().RequestId); requestId != nil {
		event.Extra["requestId"] = requestId
	}
	if platformType := ctx.Value(state.HttpHeaders().PlatformType); platformType != nil {
		event.Extra["platformType"] = platformType
	}
	if platform := ctx.Value(state.HttpHeaders().Platform); platform != nil {
		event.Extra["platform"] = platform
	}
	if version := ctx.Value(state.HttpHeaders().Version); version != nil {
		event.Extra["version"] = version
	}

	event.Message = fmt.Sprint(args...)
	hub.CaptureEvent(event)
}

// Trace logs a message at Trace level
func Trace(ctx context.Context, args ...interface{}) {
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.TraceLevel, message)
}

// Tracef logs a formatted message at Trace level
func Tracef(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.TraceLevel, message)
}

// Debug logs a message at Debug level
func Debug(ctx context.Context, args ...interface{}) {
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.DebugLevel, message)
}

// Debugf logs a formatted message at Debug level
func Debugf(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.DebugLevel, message)
}

// Info logs a message at Info level
func Info(ctx context.Context, args ...interface{}) {
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.InfoLevel, message)
}

// Infof logs a formatted message at Info level
func Infof(ctx context.Context, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.InfoLevel, message)
}

// Warn logs a message at Warn level and captures context for Sentry
func Warn(ctx context.Context, args ...interface{}) {
	// CaptureSentryEventFromContext(ctx, args...) // todo: don't use this, it will cause minimize log into sentry
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.WarnLevel, message)
}

// Warnf logs a formatted message at Warn level and captures context for Sentry
func Warnf(ctx context.Context, format string, args ...interface{}) {
	// CaptureSentryEventFromContext(ctx, args...) // todo: don't use this, it will cause minimize log into sentry
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.WarnLevel, message)
}

// Error logs a message at Error level and captures context for Sentry
func Error(ctx context.Context, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.ErrorLevel, message)
}

// Errorf logs a formatted message at Error level and captures context for Sentry
func Errorf(ctx context.Context, format string, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.ErrorLevel, message)
}

// Fatal logs a message at Fatal level and captures context for Sentry
func Fatal(ctx context.Context, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.FatalLevel, message)
}

// Fatalf logs a formatted message at Fatal level and captures context for Sentry
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	message := fmt.Sprintf(format, args...)
	logWithFields(ctx, logrus.FatalLevel, message)
}

// Panic logs a message at Panic level and captures context for Sentry
func Panic(ctx context.Context, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	message := fmt.Sprint(args...)
	logWithFields(ctx, logrus.PanicLevel, message)
}

// Panicf logs a formatted message at Panic level and captures context for Sentry
func Panicf(ctx context.Context, format string, args ...interface{}) {
	CaptureSentryEventFromContext(ctx, args...)
	logWithFields(ctx, logrus.PanicLevel, format, args...)
}

// withFields extracts fields from the context for logging
func withFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	if requestId := ctx.Value(state.HttpHeaders().RequestId); requestId != nil {
		fields[LoggerField().RequestId] = requestId
	}
	if platformType := ctx.Value(state.HttpHeaders().PlatformType); platformType != nil {
		fields[LoggerField().PlatformType] = platformType
	}
	if platform := ctx.Value(state.HttpHeaders().Platform); platform != nil {
		fields[LoggerField().Platform] = platform
	}
	if version := ctx.Value(state.HttpHeaders().Version); version != nil {
		fields[LoggerField().Version] = version
	}

	return fields
}

func funcCallerName(index int) string {
	const (
		allocatePtr = 15
		slashRune   = '/'
	)

	pc := make([]uintptr, allocatePtr)
	n := runtime.Callers(index, pc)

	frames := runtime.CallersFrames(pc[:n])
	_, _ = frames.Next()

	f := runtime.FuncForPC(pc[0])
	fName := f.Name()

	if strings.Contains(fName, string(slashRune)) {
		lastSlash := strings.LastIndexByte(fName, slashRune)
		return fName[lastSlash+1:]
	}
	return fName
}

// Modify the enableColorOutput function
func enableColorOutput() {
	Logger.SetFormatter(&ColoredJSONFormatter{
		JSONFormatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
}

type colorFunc func(string) string

// New function to log messages with color
func ColorLog(level string, message string, data map[string]interface{}) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339)
	levelColor := getColorForLevel(level) // Assume this function returns the appropriate color code for the log level

	// Format the message with colors and structure
	msg := []byte{}

	// Add timestamp
	msg = append(msg, []byte(
		"\x1b[36m"+timestamp+"\x1b[0m")..., // Cyan timestamp
	)
	msg = append(msg, []byte(" | ")...)

	// Add level
	msg = append(msg, []byte(
		"\x1b["+string(rune(levelColor))+"m"+level+"\x1b[0m")...,
	)
	msg = append(msg, []byte(" | ")...)

	// Add other fields if they exist
	if len(data) > 0 {
		msg = append(msg, []byte("\x1b[35m")...) // Magenta for fields
		for k, v := range data {
			msg = append(msg, []byte(fmt.Sprintf("%s=%v ", k, v))...)
		}
		msg = append(msg, []byte("\x1b[0m")...)
		msg = append(msg, []byte("| ")...)
	}

	// Add the message
	msg = append(msg, []byte("\x1b[37m"+message+"\x1b[0m")...) // White message
	msg = append(msg, []byte("\n")...)

	return msg, nil
}

// Helper function to get color code based on log level
func getColorForLevel(level string) int {
	switch level {
	case "DEBUG":
		return 34 // Blue
	case "INFO":
		return 32 // Green
	case "WARN":
		return 33 // Yellow
	case "ERROR":
		return 31 // Red
	default:
		return 37 // White
	}
}

// Colors for different log levels
const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorBlue   = 34
	colorGray   = 37
)

// OpenTelemetry tracing functions

// StartSpan starts an OpenTelemetry span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordSpanError records an error on the current span
func RecordSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
