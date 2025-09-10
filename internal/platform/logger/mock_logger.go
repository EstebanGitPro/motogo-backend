package logger


type MockLogger struct {
	InfoCalls  []LogCall
	ErrorCalls []LogCall
	DebugCalls []LogCall
}

type LogCall struct {
	Message string
	Args    []any
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		InfoCalls:  make([]LogCall, 0),
		ErrorCalls: make([]LogCall, 0),
		DebugCalls: make([]LogCall, 0),
	}
}

func (m *MockLogger) Info(msg string, args ...any) {
	m.InfoCalls = append(m.InfoCalls, LogCall{Message: msg, Args: args})
}

func (m *MockLogger) Error(msg string, args ...any) {
	m.ErrorCalls = append(m.ErrorCalls, LogCall{Message: msg, Args: args})
}

func (m *MockLogger) Debug(msg string, args ...any) {
	m.DebugCalls = append(m.DebugCalls, LogCall{Message: msg, Args: args})
}


func (m *MockLogger) GetLastError() *LogCall {
	if len(m.ErrorCalls) == 0 {
		return nil
	}
	return &m.ErrorCalls[len(m.ErrorCalls)-1]
}

func (m *MockLogger) GetLastInfo() *LogCall {
	if len(m.InfoCalls) == 0 {
		return nil
	}
	return &m.InfoCalls[len(m.InfoCalls)-1]
}

func (m *MockLogger) Reset() {
	m.InfoCalls = make([]LogCall, 0)
	m.ErrorCalls = make([]LogCall, 0)
	m.DebugCalls = make([]LogCall, 0)
}
