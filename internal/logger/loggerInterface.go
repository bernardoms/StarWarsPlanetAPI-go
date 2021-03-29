package logger

import "net/http"

type Interface interface {
	LogWithFields(req *http.Request, level string, fields map[string]interface{}, message string)
}
