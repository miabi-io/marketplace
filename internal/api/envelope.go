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

import "github.com/jkaninda/okapi"

// Envelope is the shared response wrapper (mirrors Miabi). It is generic so the
// concrete data shape flows into the generated OpenAPI schema — register a route
// with e.g. &Envelope[catalog.Page]{} to document the full body.
type Envelope[T any] struct {
	Success bool      `json:"success"`
	Data    T         `json:"data"`
	Error   *APIError `json:"error"`
}

// APIError is a structured error with a stable string code.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse documents a failed response (success=false, populated error) for
// OpenAPI; it is the body shape fail() emits.
type ErrorResponse = Envelope[any]

func ok(c *okapi.Context, data any) error {
	return c.JSON(200, Envelope[any]{Success: true, Data: data})
}

func fail(c *okapi.Context, status int, code, message string) error {
	return c.JSON(status, Envelope[any]{Success: false, Error: &APIError{Code: code, Message: message}})
}
