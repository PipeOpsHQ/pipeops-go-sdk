# Logging

Add logging to the SDK for debugging and monitoring.

## Enable Logging

Implement the Logger interface:

```go
import "log"

type MyLogger struct{}

func (l *MyLogger) Debug(msg string, keysAndValues ...interface{}) {
    log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Info(msg string, keysAndValues ...interface{}) {
    log.Printf("[INFO] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Warn(msg string, keysAndValues ...interface{}) {
    log.Printf("[WARN] %s %v", msg, keysAndValues)
}

func (l *MyLogger) Error(msg string, keysAndValues ...interface{}) {
    log.Printf("[ERROR] %s %v", msg, keysAndValues)
}

// Use the logger
client, _ := pipeops.NewClient("",
    pipeops.WithLogger(&MyLogger{}),
)
```

## Structured Logging

Use structured logging libraries:

```go
import "go.uber.org/zap"

type ZapLogger struct {
    logger *zap.SugaredLogger
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
    l.logger.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
    l.logger.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
    l.logger.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
    l.logger.Errorw(msg, keysAndValues...)
}

// Create logger
zapLogger, _ := zap.NewProduction()
logger := &ZapLogger{logger: zapLogger.Sugar()}

client, _ := pipeops.NewClient("",
    pipeops.WithLogger(logger),
)
```

## Log Levels

Control log verbosity:

```go
type LeveledLogger struct {
    level string
}

func (l *LeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
    if l.level == "debug" {
        log.Printf("[DEBUG] %s", msg)
    }
}

func (l *LeveledLogger) Info(msg string, keysAndValues ...interface{}) {
    if l.level == "debug" || l.level == "info" {
        log.Printf("[INFO] %s", msg)
    }
}
```

## See Also

- [Error Handling](error-handling.md)
- [Configuration](../getting-started/configuration.md)
