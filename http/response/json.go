package response

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func JSON(rw http.ResponseWriter, code int, result interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	rr := ResponseResult{}
	if code == 0 || code == http.StatusOK {
		rw.WriteHeader(http.StatusOK)
		rr.Type = ResponseResultSuccess
		if message == "" {
			message = "ok"
		}

		if result == nil {
			result = "ok"
		}
	} else {
		if statusText := http.StatusText(code); statusText != "" {
			rw.WriteHeader(code)
			if message == "" {
				message = statusText
			}
		} else {
			rw.WriteHeader(http.StatusOK)
		}

		rr.Type = ResponseResultError
	}

	rr.Code = code
	rr.Message = message
	if result != nil {
		switch res := result.(type) {
		case string:
			if u, err := url.Parse(res); err == nil && u.IsAbs() {
				rr.Url = res
			} else {
				rr.Result = res
			}
		default:
			rr.Result = res
		}
	}

	jsonData, err := json.Marshal(rr)
	if err != nil {
		rw.Write([]byte(`{"code":500,"type":"error","message":"invalid json format"}`))
		return
	}

	rw.Write(jsonData)
}
