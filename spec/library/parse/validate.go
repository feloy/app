package parse

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed app.json
var devfileSpec string

func ValidateSchema(devfile string) ([]gojsonschema.ResultError, error) {
	yamlFile, err := os.Open(devfile)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = yamlFile.Close()
	}()
	var doc map[string]interface{}
	dec := yaml.NewDecoder(yamlFile)
	err = dec.Decode(&doc)
	if err != nil {
		return nil, err
	}

	documentLoader := gojsonschema.NewGoLoader(doc)
	schemaLoader := gojsonschema.NewStringLoader(devfileSpec)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
		return nil, nil
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return result.Errors(), nil
	}
}
