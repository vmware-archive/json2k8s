package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	err := Command{}.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}
}

type Command struct{}

func (f Command) Run(args []string) error {
	var files []File

	if len(args) > 0 {
		for _, arg := range args {
			if arg == "-" {
				files = []File{File{"stdin", os.Stdin}}
			} else {
				files = append(files, File{arg, nil})
			}
		}
	} else {
		files = []File{File{"stdin", os.Stdin}}
	}

	result := V1List{
		APIVersion: "v1",
		Kind:       "List",
	}

	for _, file := range files {
		bs, err := file.Bytes()
		if err != nil {
			return err
		}

		decoder := json.NewDecoder(bytes.NewReader(bs))

		for {
			var input interface{}

			err := decoder.Decode(&input)
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("Unmarshaling file '%s': %s", file.Path, err)
			}

			result.Items = append(result.Items, f.collectResources(input)...)
		}
	}

	bs, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s", bs)

	return nil
}

func (f Command) collectResources(in interface{}) []interface{} {
	const (
		kindKey       = "kind"
		apiVersionKey = "apiVersion"
	)

	var result []interface{}

	switch typedIn := in.(type) {
	case map[string]interface{}:
		_, kindPresent := typedIn[kindKey]
		_, apiVersionPresent := typedIn[apiVersionKey]
		if kindPresent && apiVersionPresent {
			return []interface{}{typedIn}
		}
		for _, v := range typedIn {
			result = append(result, f.collectResources(v)...)
		}

	case []interface{}:
		for _, v := range typedIn {
			result = append(result, f.collectResources(v)...)
		}

	default:
		// do nothing
	}

	return result
}

type File struct {
	Path string
	out  io.Reader
}

func (f File) Bytes() ([]byte, error) {
	var out io.Reader = f.out

	if out == nil {
		actualFile, err := os.Open(f.Path)
		if err != nil {
			return nil, fmt.Errorf("Opening file '%s': %s", f.Path, err)
		}
		defer actualFile.Close()
		out = actualFile
	}

	bs, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, fmt.Errorf("Reading file '%s': %s", f.Path, err)
	}

	return bs, nil
}

type V1List struct {
	Kind       string        `json:"kind"`
	APIVersion string        `json:"apiVersion"`
	Items      []interface{} `json:"items"`
}
