package jsonparser

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

type JSONParser struct {
	JSONSchema        any
	ConvertFunction   func(any, string) any
	ValidatorFunction func(any) error
	Output            any
	Name              string
}

func (parser *JSONParser) Parse(data []byte, url string) (any, error) {
	// Check if schema and output are pointers
	if reflect.ValueOf(parser.JSONSchema).Kind() != reflect.Ptr {
		return nil, errors.New("JSONSchema must be a pointer")
	}
	if reflect.ValueOf(parser.Output).Kind() != reflect.Ptr {
		return nil, errors.New("output must be a pointer")
	}

	// Check if schema and output are structs
	if reflect.ValueOf(parser.JSONSchema).Elem().Kind() != reflect.Struct {
		return nil, errors.New("JSONSchema must be a pointer to a struct")
	}
	if reflect.ValueOf(parser.Output).Elem().Kind() != reflect.Struct {
		return nil, errors.New("output must be a pointer to a struct")
	}

	// Create a new instance of the schema and output struct
	schema := reflect.New(reflect.ValueOf(parser.JSONSchema).Elem().Type()).Interface()

	// Unmarshal the data into the schema
	err := json.Unmarshal(data, schema)
	if err != nil {
		return nil, err
	}

	// Check if the schema is valid
	err = parser.ValidatorFunction(schema)
	if err != nil {
		// Log failing url
		log.Println(url)
		return nil, err
	}

	// Use the convert function if it is set
	if parser.ConvertFunction == nil {
		return nil, errors.New("ConvertFunction is not set")
	}

	output := parser.ConvertFunction(schema, url)
	return output, nil
}
