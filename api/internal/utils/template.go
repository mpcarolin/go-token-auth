package utils

import (
	"bytes"
	"text/template"
)

func Template(tmplStr string, data map[string]interface{}) (string, error) {
	var output bytes.Buffer;

	tmpl := template.Must(template.New("parsed").Parse(tmplStr));

	err := tmpl.Execute(&output, data);
	if err != nil {
		return "", err
	}

	return output.String(), nil
}