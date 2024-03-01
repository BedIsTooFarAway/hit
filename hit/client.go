package hit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client sends http requests and returns an aggregated performance result.
// ! The fields should not be changed after initializing.
type Client struct {
	C   int // C is concurrency level
	RPS int // Throttles requests per second
}

// Do sends n http requests and returns an aggregated performance result.
func (c *Client) Do(r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(r, n)
	fmt.Printf("Finalizing with total time %d", time.Since(t))
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(r *http.Request, n int) *Result {
	p := produce(n, func() *http.Request {
		return r.Clone(context.TODO())
	})
	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
	}
	var sum Result
	pipe := split(p, c.C, Send)
	for result := range pipe {
		sum.Merge(result)
	}
	return &sum
}
