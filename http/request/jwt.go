package request

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

func GetClaims(r *http.Request) (claims *jwt.StandardClaims) {
	if tokenString, found := BearerAuth(r); found {
		if token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwt.StandardClaims{}); err == nil {
			claims = token.Claims.(*jwt.StandardClaims)
		}
	}

	return
}
