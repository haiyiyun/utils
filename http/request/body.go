package request

import (
	"bytes"
	"io"
	"net/http"
)

func GetBody(r *http.Request) ([]byte, error) {
	data, err := io.ReadAll(r.Body)
	if err == nil {
		buf := bytes.NewBuffer(data)
		r.Body = io.NopCloser(buf)
	}

	return data, err
}
