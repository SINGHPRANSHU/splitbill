package common

import "encoding/json"

type HttpErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func GetHttpErrorResponse(code int, message string) string {
	data, err :=  json.Marshal(HttpErrorResponse{Code: code, Message: message})
	if err != nil {
		return ""
	}
	return string(data)
}