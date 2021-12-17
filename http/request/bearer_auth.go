package request

import (
	"net/http"
	"strings"
)

func BearerAuth(r *http.Request, args ...string) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	if token == "" {
		queryName := "token"
		if len(args) > 0 {
			queryName = args[0]
		}

		token = r.URL.Query().Get(queryName)
	}

	return token, token != ""
}
