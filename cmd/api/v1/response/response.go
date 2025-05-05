package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Envelope map[string]any

func Write(w http.ResponseWriter, status int, content any) error {
	if nil == content {
		w.WriteHeader(status)
		return nil
	}

	data, err := json.Marshal(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	content := struct {
		Code string `json:"code"`
		Msg  string `json:"message"`
	}{
		Code: http.StatusText(status),
		Msg:  err.Error(),
	}

	return Write(w, status, Envelope{"error": content})
}
