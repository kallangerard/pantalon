package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// pathGlobs implements flag.Value for a repeatable --path-glob flag.
type pathGlobs []string

func (p *pathGlobs) String() string { return fmt.Sprintf("%v", *p) }
func (p *pathGlobs) Set(v string) error {
	*p = append(*p, v)
	return nil
}

// selectorFlag implements flag.Value for --selector/-l. It rejects blank or
// repeated values so a mis-templated CI variable cannot silently select every
// configuration.
type selectorFlag struct {
	value string
	set   bool
}

func (s *selectorFlag) String() string { return s.value }
func (s *selectorFlag) Set(v string) error {
	if s.set {
		return errors.New("selector may only be given once (use commas to combine requirements)")
	}
	if strings.TrimSpace(v) == "" {
		return errors.New("selector must not be empty; omit the flag to select all configurations")
	}
	s.value, s.set = v, true
	return nil
}

func usage(fs *flag.FlagSet) func() {
	return func() {
		w := fs.Output()
		_, _ = fmt.Fprintf(w, `pantalon - identify Terraform root module configurations for CI/CD pipelines

Walks the repository from the current directory, finds pantalon.yaml marker
files, and emits a machine-readable list of Terraform root modules suitable
for use in GitHub Actions job matrices or other CI/CD tooling.

Usage:
  pantalon [flags]

Flags:
`)
		fs.PrintDefaults()
		_, _ = fmt.Fprintf(w, `
Examples:
  pantalon
  pantalon --output-format=yaml
  pantalon --changed-dirs='["terraform/compute/environments/dev"]'
  pantalon --path-glob='terraform/compute/**'
  pantalon --path-glob='terraform/compute/**' --path-glob='terraform/data/**'
  pantalon -l 'environment=prod'
  pantalon --selector='environment in (dev,qa),team!=data'
`)
	}
}
