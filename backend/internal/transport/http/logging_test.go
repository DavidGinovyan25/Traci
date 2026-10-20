package httptransport

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type testLogs struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (l *testLogs) Write(data []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buffer.Write(data)
}

func (l *testLogs) text() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buffer.String()
}

func (l *testLogs) entries(t *testing.T) []map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewBufferString(l.text()))
	var entries []map[string]any
	for {
		var entry map[string]any
		err := decoder.Decode(&entry)
		if err == io.EOF {
			return entries
		}
		require.NoError(t, err)
		entries = append(entries, entry)
	}
}

func newTestLogger() (*slog.Logger, *testLogs) {
	logs := &testLogs{}
	return slog.New(slog.NewJSONHandler(logs, nil)), logs
}
