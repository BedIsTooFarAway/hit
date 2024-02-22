package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/BedIsTooFarAway/hit/hit"
)

const (
	bannerText = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__  _\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/
`
)

func banner() string { return bannerText[1:] }

func main() {

	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {

	d, err := time.ParseDuration("39s")
	if err != nil {
		return err
	}
	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
		t: d,
		m: "GET",
	}

	if err = f.parse(s, args); err != nil {
		return err
	}
	fmt.Fprintln(out, banner())
	fmt.Fprintf(out, "Headers: %s\n", (*headers)(&f.H))
	fmt.Fprintf(out, "Making %d %s requests to %s with a concurrency level of %d (Timeout=%s).\n", f.n, f.m, f.url, f.c, f.t)

	var sum hit.Result
	sum.Merge(&hit.Result{
		Bytes:    1000,
		Status:   http.StatusOK, // 200
		Duration: time.Second,
	})

	sum.Merge(&hit.Result{
		Bytes:    1000,
		Status:   http.StatusOK, // 200
		Duration: time.Second,
	})

	sum.Merge(&hit.Result{
		Status:   http.StatusTeapot, // 200
		Duration: 2 * time.Second,
	})

	sum.Finalize(2 * time.Second)
	//sum.Fprint(out)
	fmt.Printf("Sum is %T:\n%s", sum, sum)

	return nil
}
