package request

import (
	"net/http"
	"net/url"
)

func ParseDeleteForm(r *http.Request) (url.Values, error) {
	if data, err := GetBody(r); err == nil {
		return url.ParseQuery(string(data))
	} else {
		return nil, err
	}
}
