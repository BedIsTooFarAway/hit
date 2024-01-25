package hit

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Result is a request's stats
type Result struct {
	RPS      float64       // Requests Per Second
	Requests int           // Requests, Total
	Errors   int           // Number of errors, Total
	Bytes    int64         // Bytest Received, Total
	Duration time.Duration // Duration, Total
	Fastest  time.Duration // Request of all
	Slowest  time.Duration // Request of all
	Status   int           // HTTP status code. To be validated
	Success  float64       // Percent of non-error requests
	Error    error         // Not nil if request failed
}

// Merge this Result with another
func (r *Result) Merge(o *Result) {

	o.Requests++
	r.Bytes += o.Bytes
	switch {
	case o.Error != nil:
		fallthrough
	case o.Status > http.StatusBadRequest:
		r.Errors++
	}
	if o.Error == nil {
		r.Fastest = max(o.Duration, r.Fastest)
		r.Slowest = min(o.Duration, r.Slowest)
	}

}

func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests / int(total.Seconds()))
	r.Success = r.success()
	return r
}

func (r *Result) Fprintf(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}

	p("\nSummary:\n")
	p("\tSuccess	: %.0f%%\n", r.Success)
	p("\tRPC	: %.1f\n", r.RPS)
	p("\tRequests	: %d\n", r.Requests)
	p("\tErrors	: %d\n", r.Errors)
	p("\tBytes	: %d\n", r.Bytes)
	p("\tDuration	: %d\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest	: %d", r.Fastest)
		p("\tSlowest	: %d", r.Slowest)
	}
}

func max(a time.Duration, b time.Duration) time.Duration {
	// 0 = uninitialized
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	if a >= b {
		return a
	} else {
		return b
	}

}

func min(a time.Duration, b time.Duration) time.Duration {
	if max(a, b) == a {
		return b
	} else {
		return a
	}
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)

}

func (r *Result) success() float64 {
	req, err := float64(r.Requests), float64(r.Errors)
	return (req - err) * 100 / req
}
