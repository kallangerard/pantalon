package pantalon

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// Write renders items to w as "json" or "yaml", followed by a newline.
func Write(w io.Writer, items []ConfigurationItem, format string) error {
	var data []byte
	var err error
	switch format {
	case "json":
		data, err = yaml.MarshalWithOptions(items, yaml.JSON())
		if err != nil {
			return wrap("Error marshaling json", err)
		}
	case "yaml":
		data, err = yaml.Marshal(items)
		if err != nil {
			return wrap("Error marshaling yaml", err)
		}
	default:
		return wrap("Unsupported output format", errors.New(format))
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

// wrap prefixes err with msg. The prefixes are user-visible output, hence capitalised.
func wrap(msg string, err error) error { return fmt.Errorf("%s: %w", msg, err) }
