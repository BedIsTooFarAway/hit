package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const usageText = `
	Usage
		hit [options] url
		Options:`

// methods
const (
	get  method = "GET"
	post method = "POST"
	put  method = "PUT"
)

type flags struct {
	url  string
	n, c int
	t    time.Duration
	m    string
	H    []string
}

// number is a natural number.
// Interface: flag.Value
type number int

// method is an allowed action.
type method string

type headers []string

// toNumber is a convenience function for converting int to number.
func toNumber(p *int) *number {
	return (*number)(p)
}

func toMethod(v *string) *method {
	return (*method)(v)
}

func toHeaders(h *[]string) *headers {
	//var r headers = []string{"toHeaders"}
	return (*headers)(h)
}

// Interface Value for type headers
func (h *headers) Set(s string) error {
	*h = append(*h, s)
	return nil
}

func (h *headers) String() string {
	//fmt.Println("headers to string")

	var result []string
	for _, item := range *h {
		result = append(result, fmt.Sprintf("%q", item))
	}
	return fmt.Sprintf(strings.Join(result, ", "))
	//return string(*h)
}

//> Interface

// Interface Value for type number
func (n *number) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("should be positive")
	}
	*n = number(v)
	return err
}

func (n *number) String() string {
	return strconv.Itoa(int(*n))
}

//> Interface

// Interface Value for type method
func (v *method) Set(s string) error {
	if err := validateMethod(s); err != nil {
		return err
	}
	*v = method(strings.ToUpper(s))
	return nil
}

func (v *method) String() string {
	return (string)(*v)
}

func (f *flags) parse(s *flag.FlagSet, args []string) error {
	flag.Usage = func() {
		fmt.Fprintln(s.Output(), usageText[1:])
		s.PrintDefaults()

	}
	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
	s.DurationVar(&f.t, "t", f.t, "Timeout")
	s.Var(toMethod(&f.m), "m", "Method")
	s.Var(toHeaders(&f.H), "H", "Headers. Multiple-entry parameter")
	if err := s.Parse(args); err != nil {
		return err
	}

	f.url = s.Arg(0)

	if err := f.validate(); err != nil {
		fmt.Fprintln(s.Output(), err)
		s.Usage()
		return err
	}
	return nil
}

func (f *flags) validate() error {
	if err := validateUrl(f.url); err != nil {
		return fmt.Errorf("url: %w", err)
	}
	if err := validateMethod(f.m); err != nil {
		return fmt.Errorf("%w: %s", err, f.m)
	}

	if f.c > f.n {
		return fmt.Errorf("-c=%d: should be less than or equal to -n=%d", f.c, f.n)
	}
	return nil
}

func validateUrl(s string) error {
	u, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("required")
	case err != nil:
		err = errors.New("parse error")
	case u.Scheme != "http":
		err = errors.New(fmt.Sprintf("only supported scheme is http."))
	case u.Host == "":
		err = errors.New("missing host")
	}
	return err
}

func validateMethod(v string) error {
	methods := map[string]method{
		"GET":  get,
		"POST": post,
		"PUT":  put,
	}

	if v != "" && methods[strings.ToUpper(v)] == "" {
		return errors.New(fmt.Sprintf("Incorrect method: %s", v))
	}

	return nil
}
