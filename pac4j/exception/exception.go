// Package exception ports the pac4j exceptions referenced by the bridge:
// TechnicalException and the HttpAction control-flow exception.
package exception

import "fmt"

// TechnicalException ports org.pac4j.core.exception.TechnicalException: an
// unchecked exception wrapping a lower-level failure. In Go it is an error value.
type TechnicalException struct {
	// Message is the human-readable detail.
	Message string
	// Cause is the wrapped error, if any (Throwable cause in Java).
	Cause error
}

// NewTechnicalException wraps a cause, mirroring new TechnicalException(Throwable).
func NewTechnicalException(cause error) *TechnicalException {
	msg := ""
	if cause != nil {
		msg = cause.Error()
	}
	return &TechnicalException{Message: msg, Cause: cause}
}

// NewTechnicalExceptionMessage mirrors new TechnicalException(String).
func NewTechnicalExceptionMessage(message string) *TechnicalException {
	return &TechnicalException{Message: message}
}

// Error implements error.
func (e *TechnicalException) Error() string {
	if e.Cause != nil && e.Message == "" {
		return e.Cause.Error()
	}
	return e.Message
}

// Unwrap exposes the wrapped cause for errors.Is/As.
func (e *TechnicalException) Unwrap() error { return e.Cause }

// HttpAction ports org.pac4j.core.exception.http.HttpAction: a checked exception
// used by pac4j for control flow (e.g. redirects). populateSubject catches it and
// rewraps it as a TechnicalException.
type HttpAction struct {
	// Code is the HTTP status code carried by the action.
	Code int
	// Message is the optional detail.
	Message string
}

// Error implements error.
func (e *HttpAction) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("HttpAction %d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("HttpAction %d", e.Code)
}
