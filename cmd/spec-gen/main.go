// cmd/spec-gen converts docs/openapi.yaml to docs/openapi.json.
// Run from the module root: go run ./cmd/spec-gen
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read docs/openapi.yaml: %v\n", err)
		os.Exit(1)
	}

	var v any
	if err := yaml.Unmarshal(data, &v); err != nil {
		fmt.Fprintf(os.Stderr, "parse yaml: %v\n", err)
		os.Exit(1)
	}

	enc, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal json: %v\n", err)
		os.Exit(1)
	}
	enc = append(enc, '\n')

	if err := os.WriteFile("docs/openapi.json", enc, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "write docs/openapi.json: %v\n", err)
		os.Exit(1)
	}
}
