package logger

import "go.uber.org/zap"

// Custom logger using zap for cron using github.com/robfig/cron/v3
type CronLogger struct {
	logger *zap.Logger
}

func NewCronLogger(logger *zap.Logger) *CronLogger {
	return &CronLogger{
		logger,
	}
}

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, zap.Any("detail", keysAndValues))
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	// TODO: notify user if an error occurred

	l.logger.Error(msg,
		zap.Error(err),
		zap.Any("detail", keysAndValues),
	)
}
