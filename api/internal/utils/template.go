package utils

import (
	"bytes"
	"text/template"
)

// just a helper function to ease substituting values in a template string
func Template(tmplStr string, data map[string]any) (string, error) {
	var output bytes.Buffer;

	tmpl := template.Must(template.New("parsed").Parse(tmplStr));

	err := tmpl.Execute(&output, data);
	if err != nil {
		return "", err
	}

	return output.String(), nil
}