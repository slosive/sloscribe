package generate

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	k8syaml "sigs.k8s.io/yaml"
	"strings"
)

// DefaultServiceDefinitionDir is the default filename for the output file
const DefaultServiceDefinitionDir = "slo_definitions"

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

// WriteK8Specifications write the k8s service spec bytes to a specific writer, stdout or file
func WriteK8Specifications(writer io.Writer, header []byte, specs map[string]any, toFile bool, outputDirectory string, formats ...string) error {
	for specName, spec := range specs {
		for _, format := range formats {
			var files = make(map[string][]byte, len(formats))

			format = strings.ToLower(strings.TrimSpace(format))
			switch format {
			case "yaml":
				body, err := k8syaml.Marshal(spec)
				if err != nil {
					return err
				}

				file := filepath.Join([]string{outputDirectory, DefaultServiceDefinitionDir, fmt.Sprintf("%s.%s", specName, format)}...)
				files[file] = bytes.Join([][]byte{[]byte("---"), header, body}, []byte("\n"))
				if err := clean(file); err != nil {
					return err
				}
			default:
				return ErrUnsupportedFormat
			}

			if toFile {
				if err := writeToFile(files); err != nil {
					// @aloe code write_artefacts_error
					// @aloe title Error Creating Artefacts
					// @aloe summary The tool has failed to print outputDirectory the Sloth definitions for service.
					// @aloe details The tool has failed to print outputDirectory the Sloth definitions for service.
					return err
				}
				continue
			}

			if err := write(writer, files); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteSpecifications write the service spec bytes to a specific writer, stdout or file
func WriteSpecifications(writer io.Writer, header []byte, specs map[string]any, toFile bool, outputDirectory string, formats ...string) error {
	for specName, spec := range specs {
		for _, format := range formats {
			var files = make(map[string][]byte, len(formats))

			format = strings.ToLower(strings.TrimSpace(format))
			switch format {
			case "json":
				body, err := json.Marshal(spec)
				if err != nil {
					return err
				}
				file := filepath.Join([]string{outputDirectory, DefaultServiceDefinitionDir, fmt.Sprintf("%s.%s", specName, format)}...)
				files[file] = bytes.Join([][]byte{header, body}, []byte("\n"))
				if err := clean(file); err != nil {
					// @aloe code clean_artefacts_error
					// @aloe title Error Removing Previous Artefacts
					// @aloe summary The tool has failed to delete the artefacts from the previous execution.
					// @aloe details The tool has failed to delete the artefacts from the previous execution.
					// Try manually deleting them before running the tool again.
					return err
				}
			case "yaml":
				body, err := yaml.Marshal(spec)
				if err != nil {
					return err
				}

				file := filepath.Join([]string{outputDirectory, DefaultServiceDefinitionDir, fmt.Sprintf("%s.%s", specName, format)}...)
				files[file] = bytes.Join([][]byte{[]byte("---"), header, body}, []byte("\n"))
				if err := clean(file); err != nil {
					return err
				}
			default:
				return ErrUnsupportedFormat
			}

			if toFile {
				if err := writeToFile(files); err != nil {
					// @aloe code write_artefacts_error
					// @aloe title Error Creating Artefacts
					// @aloe summary The tool has failed to print outputDirectory the Sloth definitions for service.
					// @aloe details The tool has failed to print outputDirectory the Sloth definitions for service.
					return err
				}
				continue
			}

			if err := write(writer, files); err != nil {
				return err
			}
		}
	}

	return nil
}

func clean(files ...string) error {
	for _, file := range files {
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

// write the files to the writer, the caller is in charge of closing the writer
func write(w io.Writer, files map[string][]byte) error {
	for _, body := range files {
		var err error
		// write to writer, this must be closed by the caller
		_, err = w.Write(body)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeToFile writes the files to the specified file paths. The function handles its writers
func writeToFile(files map[string][]byte) error {
	for path, body := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		w, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		_, err = w.Write(body)
		if err != nil {
			return err
		}
	}
	return nil
}
