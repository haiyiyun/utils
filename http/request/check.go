package request

import (
	"net/http"
	"strings"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/geo"
)

func CheckUserGuest(r *http.Request, guest bool) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		return claims.Guest == guest
	}
}

func CheckUserLevel(r *http.Request, level int) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		return claims.Level == level
	}
}

func CheckUserIP(r *http.Request, ip string) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		return claims.IP == ip
	}
}

func CheckUserUserAgentContains(r *http.Request, userAgent string) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		return strings.Contains(claims.UserAgent, userAgent)
	}
}

func CheckUserLocation(r *http.Request, location geometry.PointCoordinates, maxDistance, minDistance float64) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if claims.Location.Coordinates != geometry.NilPointCoordinates {
			meters := geo.HaversineDistance(claims.Location.Coordinates, location)
			if meters >= minDistance && minDistance <= maxDistance {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

func CheckUserRole(r *http.Request, role string) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if roles := claims.Roles; len(roles) == 0 {
			return false
		} else {
			found := false
			for _, urole := range roles {
				if urole.Role == role {
					found = true
					break
				}
			}

			return found
		}
	}
}

func CheckUserRoleAndRoleLevel(r *http.Request, role string, level int) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if roles := claims.Roles; len(roles) == 0 {
			return false
		} else {
			found := false
			for _, urole := range roles {
				if urole.Role == role && urole.Level == level {
					found = true
					break
				}
			}

			return found
		}
	}
}

func CheckUserTag(r *http.Request, tag string) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if tags := claims.Tags; len(tags) == 0 {
			return false
		} else {
			found := false
			for _, utag := range tags {
				if utag.Tag == tag {
					found = true
					break
				}
			}

			return found
		}
	}
}

func CheckUserTagAndTagLevel(r *http.Request, tag string, level int) bool {
	if claims := GetClaims(r); claims == nil {
		return false
	} else {
		if tags := claims.Tags; len(tags) == 0 {
			return false
		} else {
			found := false
			for _, utag := range tags {
				if utag.Tag == tag && utag.Level == level {
					found = true
					break
				}
			}

			return found
		}
	}
}
