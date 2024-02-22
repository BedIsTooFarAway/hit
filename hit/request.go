package hit

import (
	"fmt"
	"net/http"
	"time"
)

// Send http request and return performance stats
func Send(r *http.Request) *Result {
	t := time.Now()

	fmt.Printf("request: %s\n", r.URL)
	time.Sleep(100 * time.Millisecond)
	return &Result{
		Duration: time.Since(t),
		Bytes:    10,
		Status:   http.StatusOK,
	}
}
