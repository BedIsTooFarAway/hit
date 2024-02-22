package hit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

type ResultProt struct {
	RPS float64 // Requests Per Second
}

func (r ResultProt) String() string {
	return fmt.Sprintf("Requests per second: %f", r.RPS)
}

// Merge this Result with another
func (r *Result) Merge(o *Result) {

	r.Requests++
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

// Wrap up and aggregate
func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests / int(total.Seconds()))
	return r
}

func (r *Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}

	p("\nSummary:\n")
	p("\tSuccess\t\t: %.0f%%\n", r.success())
	p("\tRPS\t\t: %.1f\n", r.RPS)
	p("\tRequests\t: %d\n", r.Requests)
	p("\tErrors\t\t: %d\n", r.Errors)
	p("\tBytes\t\t: %d\n", r.Bytes)
	p("\tDuration\t: %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastest\t\t: %s\n", r.Fastest)
		p("\tSlowest\t\t: %s\n", r.Slowest)
	}
}

func (r Result) String() string {
	var s strings.Builder
	r.Fprint(&s)
	return s.String()
}

func (r *Result) success() float64 {
	req, err := float64(r.Requests), float64(r.Errors)
	return (req - err) * 100 / req
}

// Helpers
func min(a time.Duration, b time.Duration) time.Duration {
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

func max(a time.Duration, b time.Duration) time.Duration {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}

	if a < b {
		return a
	} else {
		return b
	}
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)

}
