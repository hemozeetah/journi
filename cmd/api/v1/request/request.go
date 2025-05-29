package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func ParseFile(r *http.Request, name string) ([]string, error) {
	res := []string{}

	files := r.MultipartForm.File[name]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileName := uuid.NewString() + filepath.Ext(fileHeader.Filename)
		filePath := filepath.Join("./cmd/api/v1/static", fileName)

		dest, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer dest.Close()

		_, err = io.Copy(dest, file)
		if err != nil {
			return nil, err
		}

		res = append(res, fileName)
	}

	return res, nil
}

func ParseForm(r *http.Request, dest any) error {
	return json.Unmarshal([]byte(r.FormValue("data")), dest)
}
