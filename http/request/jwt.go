package request

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/haiyiyun/plugins/user/predefined"
)

func GetClaims(r *http.Request) (claims *predefined.JWTTokenClaims) {
	if tokenString, found := BearerAuth(r); found {
		if token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &predefined.JWTTokenClaims{}); err == nil {
			claims = token.Claims.(*predefined.JWTTokenClaims)
		}
	}

	return
}

func GetStandardClaims(r *http.Request) (claims *jwt.StandardClaims) {
	if tokenString, found := BearerAuth(r); found {
		if token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwt.StandardClaims{}); err == nil {
			claims = token.Claims.(*jwt.StandardClaims)
		}
	}

	return
}
