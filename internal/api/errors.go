/*
 * Copyright 2026 Jonas Kaninda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package api

import (
	"errors"
	"net/http"

	"github.com/jkaninda/okapi"
)

// ErrorHandler returns an okapi.ErrorHandler that renders every framework error
// (request binding/validation failures, 404s, panics) as the same
// {success,data,error} envelope the handlers emit, so the API never leaks
// Okapi's default error shape.
func ErrorHandler() okapi.ErrorHandler {
	return func(c *okapi.Context, code int, message string, err error) error {
		var errStr string
		errCode := statusCode(code)
		if err != nil {
			errStr = err.Error()
			// An error carrying a stable machine code overrides the status-derived one.
			var coder interface{ Code() string }
			if errors.As(err, &coder) {
				if cc := coder.Code(); cc != "" {
					errCode = cc
				}
			}
			// Promote the real reason to the user-facing message when the framework
			// only handed us the generic status text.
			if message == "" || message == http.StatusText(code) {
				message = errStr
			}
		}
		if message == "" {
			message = http.StatusText(code)
		}
		return c.JSON(code, Envelope[any]{
			Success: false,
			Error:   &APIError{Code: errCode, Message: message, Error: errStr},
		})
	}
}

func statusCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case http.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	default:
		return "ERROR"
	}
}
