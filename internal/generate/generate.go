package generate

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"gopkg.in/yaml.v3"
)

// defaultOutputFile is the default filename for the output file
const defaultOutputFile = "sloth_definitions"

// ErrUnsupportedFormat is returned if the output format is unsupported
var ErrUnsupportedFormat = errors.New("the specification is in an invalid format")

func IsValidOutputFormat(format string) bool {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "json", "yaml":
		return true
	}
	return false
}

// WriteSpecification write the application bytes to a given writer
func WriteSpecification(spec *sloth.Spec, stdout bool, out string, formats ...string) error {
	// remove all previous output files
	err := cleanAll(spec.Service, formats...)
	if err != nil {
		// @aloe code clean_artifacts_error
		// @aloe title Error Removing Previous Artifacts
		// @aloe summary The tool has failed to delete the artifacts from the previous execution.
		// @aloe details The tool has failed to delete the artifacts from the previous execution.
		// Try manually deleting them before running the tool again.
		return err
	}

	outputFileName := defaultOutputFile
	for _, format := range formats {
		var files = make(map[string][]byte)
		var body []byte
		var err error

		format = strings.ToLower(strings.TrimSpace(format))
		switch format {
		case "json":
			body, err = json.Marshal(spec)
			file := filepath.Join([]string{out, fmt.Sprintf("%s.%s", outputFileName, format)}...)
			files[file] = body
		case "yaml":
			body, err = yaml.Marshal(spec)
			file := filepath.Join([]string{out, fmt.Sprintf("%s.%s", outputFileName, format)}...)
			files[file] = body
		default:
			return ErrUnsupportedFormat
		}
		if err != nil {
			return err
		}

		if err = writeAll(stdout, files); err != nil {
			// @aloe code write_artifacts_error
			// @aloe title Error Creating Artifacts
			// @aloe summary The tool has failed to print out the Sloth definitions for service.
			// @aloe details The tool has failed to print out the Sloth definitions for service.
			return err
		}
	}

	return nil
}

func cleanAll(applicationName string, formats ...string) error {
	for _, format := range formats {
		if !IsValidOutputFormat(format) {
			continue
		}

		file := fmt.Sprintf("%s.%s", defaultOutputFile, format)
		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			// delete spec file
			err = os.RemoveAll(file)
			if err != nil {
				return errors.Annotatef(err, "could not delete existing file %q", file)
			}
		}
	}
	return nil
}

func writeAll(stdout bool, files map[string][]byte) error {
	for path, body := range files {
		var err error
		var w = io.WriteCloser(os.Stdout)

		dirpath := filepath.Dir(path)
		if err := os.MkdirAll(dirpath, 0755); err != nil {
			return err
		}

		// decide which writer to use to print the application spec
		if !stdout {
			w, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		}
		// write to writer
		_, err = w.Write(body)
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}
