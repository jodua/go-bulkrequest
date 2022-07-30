package schemas

import (
	"errors"
	jsonparser "github.com/jodua/go-bulkrequest/JSONParser"
	"reflect"
)

type JSONPlaceholderTodoSchema struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type JSONPlaceholderTodo struct {
	UserId    int
	Id        int
	Title     string
	Completed bool
}

func JSONPlaceholderConvertFunction(input any, url string) any {
	output := JSONPlaceholderTodo{}

	schema := reflect.ValueOf(input).Elem()
	output.UserId = int(schema.FieldByName("UserId").Int())
	output.Id = int(schema.FieldByName("Id").Int())
	output.Title = schema.FieldByName("Title").String()
	output.Completed = schema.FieldByName("Completed").Bool()
	return output
}

func JSONPlaceholderValidatorFunction(input any) error {
	if reflect.ValueOf(input).Elem().Kind() != reflect.Struct {
		return errors.New("input must be a struct")
	}
	if reflect.ValueOf(input).Elem().Type() != reflect.TypeOf(JSONPlaceholderTodoSchema{}) {
		return errors.New("input must be a JSONPlaceholderTodoSchema")
	}

	// Check if userId is empty
	if reflect.ValueOf(input).Elem().FieldByName("UserId") == reflect.ValueOf(0) {
		return errors.New("response data is empty")
	}
	// Check if id is empty
	if reflect.ValueOf(input).Elem().FieldByName("Id") == reflect.ValueOf(0) {
		return errors.New("response data is empty")
	}
	// Check if title is empty
	if reflect.ValueOf(input).Elem().FieldByName("Title").Len() == 0 {
		return errors.New("response data is empty")
	}
	// Check if completed is empty
	if reflect.ValueOf(input).Elem().FieldByName("Completed") == reflect.ValueOf(false) {
		return errors.New("response data is empty")
	}
	return nil
}

var JSONPlaceholderTodoParser = jsonparser.JSONParser{
	ConvertFunction:   JSONPlaceholderConvertFunction,
	ValidatorFunction: JSONPlaceholderValidatorFunction,
	JSONSchema:        &JSONPlaceholderTodoSchema{},
	Output:            &JSONPlaceholderTodo{},
	Name:              "JSONPlaceholderTodoParser",
}
