package request

import (
	"net/http"

	"github.com/haiyiyun/utils/help"
)

func CheckUserRole(r *http.Request, role string) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if roles := claims.Roles; len(roles) == 0 {
			return false
		} else {
			return help.NewSlice(roles).CheckItem(role)
		}
	}
}
