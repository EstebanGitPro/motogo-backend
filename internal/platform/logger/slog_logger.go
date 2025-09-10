package logger

import "log/slog"


type SlogLogger struct{}


func NewSlogLogger() Logger {
	return &SlogLogger{}
}

func (s *SlogLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}
