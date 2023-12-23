package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type flags struct {
	url  string
	n, c int
}

// parseFunc is a command-line flag parser function
type parseFunc func(string) error

func (f *flags) parse() (err error) {
	// a map of flag names and parsers.
	parsers := map[string]parseFunc{
		"url": f.urlVar(&f.url), // parses the url flag and updates f.url
		"n":   f.intVar(&f.n),   // parses the n (number of requests) flag and updates f.n
		"c":   f.intVar(&f.c),   // parses the c (concurrent requests) flag and updates f.c
	}
	for _, arg := range os.Args[1:] {
		n, v, ok := strings.Cut(arg, "=")
		if !ok {
			continue // Can't parse the flag
		}
		parse, ok := parsers[strings.TrimPrefix(n, "-")]
		if !ok {
			continue // Can't find a parser
		}
		if err := parse(v); err != nil {
			err = fmt.Errorf("Invalid value %q for flag %s: %w", v, n, err)
			break // Parsing error
		}
	}
	return err
}

func (f *flags) urlVar(p *string) parseFunc {
	return func(s string) error {
		_, err := url.Parse(s)
		*p = s
		return err
	}
}

func (f *flags) intVar(p *int) parseFunc {
	return func(s string) (err error) {
		*p, err = strconv.Atoi(s)
		return err
	}
}
