package util

import (
	"cloudcanal-openapi-cli/internal/openapi"
	"errors"
	"strconv"
	"strings"
)

func SummarizeError(err error) string {
	if err == nil {
		return ""
	}

	var serverErr *openapi.ServerError
	if errors.As(err, &serverErr) {
		body := strings.TrimSpace(serverErr.ResponseBody)
		if body == "" {
			body = "server error"
		}
		return "HTTP " + strconv.Itoa(serverErr.StatusCode) + ": " + body
	}

	root := err
	for {
		next := errors.Unwrap(root)
		if next == nil {
			break
		}
		root = next
	}
	message := strings.TrimSpace(root.Error())
	if message == "" {
		message = strings.TrimSpace(err.Error())
	}
	if message == "" {
		return "unknown error"
	}
	return message
}
