package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data any) error {
	return validate.Struct(data)
}

func ParseBody(r *http.Request, dest any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request: unable to read payload: %w", err)
	}

	return json.Unmarshal(data, dest)
}

func ParseQueryParams(r *http.Request, dest any) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("output must be a pointer to struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if field.Type.Kind() != reflect.String {
			return errors.New("fields must be a string")
		}

		name := field.Tag.Get("param")
		if name == "" {
			continue
		}

		value := r.URL.Query().Get(name)
		fieldVal.Set(reflect.ValueOf(value))
	}

	return nil
}
