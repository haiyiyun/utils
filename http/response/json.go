package response

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/haiyiyun/log"
)

func jsonBytes(code int, result interface{}, message string) (int, []byte) {
	rr := ResponseResult{}
	var statusCode int
	if code == 0 || code == http.StatusOK {
		statusCode = http.StatusOK
		rr.Type = ResponseResultSuccess
		if message == "" {
			message = "ok"
		}

		if result == nil {
			result = "ok"
		}
	} else {
		if statusText := http.StatusText(code); statusText != "" {
			statusCode = code
			if message == "" {
				message = statusText
			}
		} else {
			statusCode = http.StatusOK
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
		log.Error(err)
		statusCode = http.StatusInternalServerError
		jsonData = []byte(`{"code":500,"type":"error","message":"invalid json format"}`)
	}

	return statusCode, jsonData
}

func JSON(rw http.ResponseWriter, code int, result interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	statusCode, jsonData := jsonBytes(code, result, message)
	rw.WriteHeader(statusCode)
	rw.Write(jsonData)
}

func JSONString(code int, result interface{}, message string) string {
	_, jsonData := jsonBytes(code, result, message)
	return string(jsonData)
}
