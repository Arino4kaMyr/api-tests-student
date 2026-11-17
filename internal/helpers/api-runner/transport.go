package api_runner

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// minRequestInterval минимальная пауза между запросами одного процесса: стенд отвечает 503 при слишком частых запросах
	minRequestInterval = 120 * time.Millisecond
	maxRetries         = 6
	retryBackoff       = 300 * time.Millisecond
)

// throttledTransport ограничивает частоту запросов и повторяет запрос, если стенд ответил 503 (rate limit на nginx)
type throttledTransport struct {
	next http.RoundTripper
}

var (
	throttleMu      sync.Mutex
	lastRequest     time.Time
	sharedTransport = &throttledTransport{next: http.DefaultTransport}
)

// RoundTrip выполняет запрос с соблюдением паузы между запросами и повторами при 503
func (tr *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}

	for attempt := 0; ; attempt++ {
		waitForSlot()
		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err := tr.next.RoundTrip(req)
		if err != nil || resp.StatusCode != http.StatusServiceUnavailable || attempt >= maxRetries {
			return resp, err
		}

		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		time.Sleep(retryBackoff * time.Duration(attempt+1))
	}
}

// waitForSlot блокирует выполнение до тех пор, пока с прошлого запроса не пройдёт minRequestInterval
func waitForSlot() {
	throttleMu.Lock()
	defer throttleMu.Unlock()

	if wait := minRequestInterval - time.Since(lastRequest); wait > 0 {
		time.Sleep(wait)
	}
	lastRequest = time.Now()
}

// readRequestBody вычитывает тело запроса, чтобы его можно было отправить повторно
func readRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}
	defer req.Body.Close()

	return io.ReadAll(req.Body)
}
