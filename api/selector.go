package api

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type operator int

const (
	opEquals operator = iota
	opNotEquals
	opIn
	opNotIn
	opExists
	opNotExists
)

type requirement struct {
	key    string
	op     operator
	values []string
}

// Selector is a parsed label selector. All requirements must match.
//
// The syntax is the set-based and equality-based subset of Kubernetes label
// selectors, as used by kubectl -l:
// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
type Selector []requirement

var setRequirementRegexp = regexp.MustCompile(`^(\S+)\s+(in|notin)\s*\((.*)\)$`)

// ParseSelector parses a selector such as "env=prod,team in (a,b),!legacy".
// An empty string yields a selector that matches everything.
func ParseSelector(s string) (Selector, error) {
	if strings.TrimSpace(s) == "" {
		return Selector{}, nil
	}

	parts, err := splitRequirements(s)
	if err != nil {
		return nil, err
	}

	sel := make(Selector, 0, len(parts))
	for _, part := range parts {
		req, err := parseRequirement(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid selector %q: %w", s, err)
		}
		sel = append(sel, req)
	}

	return sel, nil
}

// splitRequirements splits on commas that are not inside parentheses.
func splitRequirements(s string) ([]string, error) {
	var parts []string
	depth, start := 0, 0

	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
			if depth < 0 {
				return nil, fmt.Errorf("invalid selector %q: unbalanced parentheses", s)
			}
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	if depth != 0 {
		return nil, fmt.Errorf("invalid selector %q: unbalanced parentheses", s)
	}

	return append(parts, s[start:]), nil
}

func parseRequirement(s string) (requirement, error) {
	if s == "" {
		return requirement{}, fmt.Errorf("empty requirement")
	}

	if m := setRequirementRegexp.FindStringSubmatch(s); m != nil {
		op := opIn
		if m[2] == "notin" {
			op = opNotIn
		}
		if strings.TrimSpace(m[3]) == "" {
			return requirement{}, fmt.Errorf("%s requires a non-empty set of values", m[2])
		}
		return newRequirement(m[1], op, strings.Split(m[3], ","))
	}

	if i := strings.Index(s, "!="); i >= 0 {
		return newRequirement(s[:i], opNotEquals, []string{s[i+2:]})
	}
	if i := strings.Index(s, "=="); i >= 0 {
		return newRequirement(s[:i], opEquals, []string{s[i+2:]})
	}
	if i := strings.Index(s, "="); i >= 0 {
		return newRequirement(s[:i], opEquals, []string{s[i+1:]})
	}

	if key, found := strings.CutPrefix(s, "!"); found {
		return newRequirement(key, opNotExists, nil)
	}

	return newRequirement(s, opExists, nil)
}

func newRequirement(key string, op operator, values []string) (requirement, error) {
	key = strings.TrimSpace(key)
	if !isValidLabelKey(key) {
		return requirement{}, fmt.Errorf("invalid label key %q", key)
	}

	for i, v := range values {
		v = strings.TrimSpace(v)
		if !isValidLabelValue(v) {
			return requirement{}, fmt.Errorf("invalid label value %q", v)
		}
		values[i] = v
	}

	return requirement{key: key, op: op, values: values}, nil
}

// Matches reports whether the labels satisfy every requirement.
func (s Selector) Matches(labels map[string]string) bool {
	for _, req := range s {
		if !req.matches(labels) {
			return false
		}
	}
	return true
}

func (r requirement) matches(labels map[string]string) bool {
	v, ok := labels[r.key]

	switch r.op {
	case opEquals:
		return ok && v == r.values[0]
	case opNotEquals:
		return !ok || v != r.values[0]
	case opIn:
		return ok && slices.Contains(r.values, v)
	case opNotIn:
		return !ok || !slices.Contains(r.values, v)
	case opExists:
		return ok
	case opNotExists:
		return !ok
	}
	return false
}
