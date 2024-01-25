package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
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
		//H: []string{"sss"},
	}

	if err = f.parse(s, args); err != nil {
		return err
	}
	fmt.Fprintln(out, banner())
	fmt.Fprintf(out, "Headers: %s\n", (*headers)(&f.H))
	fmt.Fprintf(out, "Making %d %s requests to %s with a concurrency level of %d (Timeout=%s).\n", f.n, f.m, f.url, f.c, f.t)

	return nil
}
