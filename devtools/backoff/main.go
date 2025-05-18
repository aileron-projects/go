package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/aileron-projects/go/ztime/zbackoff"
	"gopkg.in/yaml.v3"
)

var (
	attempts = flag.Int("attempts", 100, "final number of attempts")
	output   = flag.String("output", "backoff.csv", "output file path")
	file     = flag.String("file", "config.yaml", "config file path")
)

func main() {
	flag.Parse()
	if *output == "" {
		flag.Usage()
		os.Exit(1)
	}
	b, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}
	config := &config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		panic(err)
	}

	durations := map[string][]string{
		"Fixed":            config.FixedBackoff.backoff(*attempts),
		"Random":           config.RandomBackoff.backoff(*attempts),
		"Linear":           config.LinearBackoff.backoff(*attempts),
		"Polynomial":       config.PolynomialBackoff.backoff(*attempts),
		"Exponential":      config.ExponentialBackoff.backoff(*attempts),
		"ExponentialFull":  config.ExponentialBackoffFullJitter.backoff(*attempts),
		"ExponentialEqual": config.ExponentialBackoffEqualJitter.backoff(*attempts),
		"Fibonacci":        config.FibonacciBackoff.backoff(*attempts),
	}

	out, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	cw := csv.NewWriter(out)
	defer cw.Flush()

	err = cw.Write([]string{
		"Fixed",
		"Random",
		"Linear",
		"Polynomial",
		"Exponential",
		"ExponentialFull",
		"ExponentialEqual",
		"Fibonacci",
	})
	if err != nil {
		panic(err)
	}
	for i := range *attempts + 1 {
		cw.Write([]string{
			durations["Fixed"][i],
			durations["Random"][i],
			durations["Linear"][i],
			durations["Polynomial"][i],
			durations["Exponential"][i],
			durations["ExponentialFull"][i],
			durations["ExponentialEqual"][i],
			durations["Fibonacci"][i],
		})
	}
}

type config struct {
	FixedBackoff                  *fixedBackoff                  `yaml:"fixedBackoff"`
	RandomBackoff                 *randomBackoff                 `yaml:"randomBackoff"`
	LinearBackoff                 *linearBackoff                 `yaml:"linearBackoff"`
	PolynomialBackoff             *polynomialBackoff             `yaml:"polynomialBackoff"`
	ExponentialBackoff            *exponentialBackoff            `yaml:"exponentialBackoff"`
	ExponentialBackoffFullJitter  *exponentialBackoffFullJitter  `yaml:"exponentialBackoffFullJitter"`
	ExponentialBackoffEqualJitter *exponentialBackoffEqualJitter `yaml:"exponentialBackoffEqualJitter"`
	FibonacciBackoff              *fibonacciBackoff              `yaml:"fibonacciBackoff"`
}

type fixedBackoff struct {
	Value int64
}

func (b *fixedBackoff) backoff(max int) []string {
	v := time.Duration(b.Value) * time.Microsecond
	backoff := zbackoff.NewFixedBackoff(v)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type randomBackoff struct {
	Offset      int64
	Fluctuation int64
}

func (b *randomBackoff) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	backoff := zbackoff.NewRandomBackoff(o, f)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type linearBackoff struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
}

func (b *linearBackoff) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Microsecond
	backoff := zbackoff.NewLinearBackoff(o, f, c)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type polynomialBackoff struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
	Exponent    float64
}

func (b *polynomialBackoff) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Microsecond
	backoff := zbackoff.NewPolynomialBackoff(o, f, c, b.Exponent)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type exponentialBackoff struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
}

func (b *exponentialBackoff) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Nanosecond
	backoff := zbackoff.NewExponentialBackoff(o, f, c)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type exponentialBackoffFullJitter struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
}

func (b *exponentialBackoffFullJitter) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Nanosecond
	backoff := zbackoff.NewExponentialBackoffFullJitter(o, f, c)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type exponentialBackoffEqualJitter struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
}

func (b *exponentialBackoffEqualJitter) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Nanosecond
	backoff := zbackoff.NewExponentialBackoffEqualJitter(o, f, c)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}

type fibonacciBackoff struct {
	Offset      int64
	Fluctuation int64
	Coeff       int64
}

func (b *fibonacciBackoff) backoff(max int) []string {
	o := time.Duration(b.Offset) * time.Microsecond
	f := time.Duration(b.Fluctuation) * time.Microsecond
	c := time.Duration(b.Coeff) * time.Nanosecond
	backoff := zbackoff.NewFibonacciBackoff(o, f, c)
	ds := make([]string, max+1)
	for i := range max + 1 {
		ds[i] = strconv.FormatInt(int64(backoff.Attempt(i)), 10)
	}
	return ds
}
