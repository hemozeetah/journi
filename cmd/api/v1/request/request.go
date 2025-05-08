package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
