package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() *Response {
	return &Response{
		Status: StatusOk,
	}
}

func Error(msg string) *Response {
	return &Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) *Response {
	var errMsg []string

	for _, val := range errs {
		switch val.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", val.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not a valid URL", val.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", val.Field()))
		}
	}

	return &Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
}
