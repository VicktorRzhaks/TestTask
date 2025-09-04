package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: jsonf field1 field2 ...")
		os.Exit(1)
	}

	fieldOrder := os.Args[1:]
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read stdin:", err)
		os.Exit(1)
	}

	var probe any

	if err := json.Unmarshal(input, &probe); err != nil {
		fmt.Fprintln(os.Stderr, "json unmarshal:", err)
		os.Exit(1)
	}

	if _, ok := probe.([]any); !ok {
		fmt.Fprintln(os.Stderr, "input must be a JSON array of objects")
		os.Exit(1)
	}
	var objects []map[string]any

	if err := json.Unmarshal(input, &objects); err != nil {
		fmt.Fprintln(os.Stderr, "json unmarshal objects:", err)
		os.Exit(1)
	}

	filtered := make([]map[string]any, len(objects))

	for i, obj := range objects {
		if obj == nil {
			fmt.Fprintln(os.Stderr, "array must contain only JSON objects")
			os.Exit(1)
		}

		out := make(map[string]any)

		for _, k := range fieldOrder {
			if v, ok := obj[k]; ok {
				out[k] = v
			}
		}
		filtered[i] = out
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
  
	if err := enc.Encode(filtered); err != nil {
		fmt.Fprintln(os.Stderr, "json encode:", err)
		os.Exit(1)
	}
}
