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
		fmt.Fprint(os.Stderr, "error occured", err)
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
	if f.rps > 0 {
		fmt.Fprintf(out, "(RPS: %d)\n", f.rps)
	}

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	c := hit.Client{
		RPS: f.rps,
		C:   f.c,
	}
	sum := c.Do(request, f.n)

	sum.Fprint(out)

	return nil
}
