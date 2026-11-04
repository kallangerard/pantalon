package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kallangerard/pantalon/internal/pantalon"
)

// errUsage means the flag package already reported the problem.
var errUsage = errors.New("usage")

func main() {
	switch err := run(os.Args[1:], os.Stdout, os.Stderr); {
	case err == nil:
	case errors.Is(err, errUsage):
		os.Exit(2)
	default:
		log.Fatal(err)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("pantalon", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = usage(fs)
	help := fs.Bool("help", false, "Show help")
	outputFormat := fs.String("output-format", "json", "Output format: json or yaml")
	changedDirs := fs.String("changed-dirs", "", `JSON array of changed directories; filters output to matching configs (e.g. '["terraform/compute/environments/dev"]')`)
	var globs pathGlobs
	fs.Var(&globs, "path-glob", "Doublestar glob pattern to filter configurations by directory path (repeatable, OR logic)")
	var selector selectorFlag
	fs.Var(&selector, "selector", "Label selector to filter configurations (e.g. 'environment=prod,team in (a,b),!legacy')")
	fs.Var(&selector, "l", "Shorthand for --selector")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil // usage already printed
		}
		return errUsage
	}

	if *help {
		fs.Usage()
		return nil
	}

	configurations, warnings, err := pantalon.Discover(os.DirFS("."))
	printWarnings(stderr, warnings)
	if err != nil {
		return wrap("Error listing configurations", err)
	}

	unfiltered, err := pantalon.MarshalItems(configurations)
	if err != nil {
		return wrap("Error marshaling items", err)
	}

	items, err := pantalon.Filter(unfiltered, pantalon.FilterOptions{
		ChangedDirs: *changedDirs,
		PathGlobs:   globs,
		Selector:    selector.value,
	})
	if err != nil {
		return wrap("Error filtering items", err)
	}

	return pantalon.Write(stdout, items, *outputFormat)
}

// printWarnings writes warnings as GitHub Actions annotations when running
// there, and as log lines otherwise.
func printWarnings(w io.Writer, warnings []pantalon.Warning) {
	logger := log.New(w, "", log.LstdFlags)
	for _, wn := range warnings {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			// Workflow commands must start the line, so bypass the log prefix.
			_, _ = fmt.Fprintf(w, "::warning file=%s,line=%d::%s\n", wn.Path, wn.Line, wn.Message)
		} else {
			logger.Printf("warning: %s:%d: %s", wn.Path, wn.Line, wn.Message)
		}
	}
}

// wrap prefixes err with msg. The prefixes are user-visible output, hence capitalised.
func wrap(msg string, err error) error { return fmt.Errorf("%s: %w", msg, err) }
