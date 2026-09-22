package common

import (
	"fmt"
	"os"
	"strings"
)

func Err(msg interface{}) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

func ErrCode(code, msg string) {
	ErrCodeWithRequest(code, msg, "")
}

func ErrCodeWithRequest(code, msg string, requestID string) {
	if code == "" {
		Err(msg)
	}
	parts := []string{"Error Code: " + code, "Description: " + ErrorCodeSummary(code)}
	if msg != "" {
		parts = append(parts, "Detail: "+msg)
	}
	if requestID != "" {
		parts = append(parts, "Request ID: "+requestID)
	}
	Err(strings.Join(parts, "\n"))
}

func ErrorCodeSummary(code string) string {
	switch code {
	case ErrorCodeInvalidRequest:
		return "request format is invalid. Check CLI arguments and generated request data."
	case ErrorCodeUnauthorized:
		return "login has expired or the session token is invalid. Run `beaker login` again."
	case ErrorCodeForbidden:
		return "the current user is not allowed to perform this operation."
	case ErrorCodeNotFound:
		return "the requested resource was not found. Check the id, alias, or URL."
	case ErrorCodeDatabase:
		return "database operation failed. Check MySQL address, credentials, network, connection limits, and whether MySQL is out of disk or memory."
	case ErrorCodeCache:
		return "cache operation failed. Check Redis address, network, memory policy, and available memory."
	case ErrorCodeInternal:
		return "server internal error. Use the request id to locate server logs."
	case ErrorCodeDecode:
		return "request encryption or decoding failed. Log in again and retry."
	case ErrorCodeUpload:
		return "image upload was rejected. Check filename, file type, and size limit."
	case ErrorCodeNetwork:
		return "CLI could not reach the server. Check URL, DNS, firewall, and network connectivity."
	case ErrorCodeHTTP:
		return "server returned an HTTP error. Check the server status and reverse proxy configuration."
	case ErrorCodeInvalidResponse:
		return "server response is not valid Beaker JSON. Check whether the URL points to the admin service."
	default:
		return "unknown error. Use the detail and request id to locate server logs."
	}
}

func Assert(e error) {
	if e != nil {
		Err(e)
	}
}
