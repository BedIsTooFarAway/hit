package hit

import (
	"net/http"
	"time"
)

// SendFunc is the type of the function called by Client.Do
// to send an HTTP request and return a performance result.
type SendFunc func(*http.Request) *Result

// Send http request and return performance stats
func Send(r *http.Request) *Result {
	t := time.Now()

	time.Sleep(100 * time.Millisecond)
	return &Result{
		Duration: time.Since(t),
		Bytes:    10,
		Status:   http.StatusOK,
	}
}
