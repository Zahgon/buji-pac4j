package exception

import (
	"errors"
	"testing"
)

// TestTechnicalExceptionFromCause covers NewTechnicalException, Error (cause
// path) and Unwrap.
func TestTechnicalExceptionFromCause(t *testing.T) {
	cause := errors.New("boom")
	te := NewTechnicalException(cause)
	if te.Error() != "boom" {
		t.Fatalf("Error() should surface cause message, got %q", te.Error())
	}
	if !errors.Is(te, cause) {
		t.Fatalf("Unwrap must expose the cause for errors.Is")
	}
	if NewTechnicalException(nil).Error() != "" {
		t.Fatalf("nil cause yields empty message")
	}
}

// TestTechnicalExceptionFromMessage covers NewTechnicalExceptionMessage and the
// message-only Error path.
func TestTechnicalExceptionFromMessage(t *testing.T) {
	te := NewTechnicalExceptionMessage("bad config")
	if te.Error() != "bad config" {
		t.Fatalf("Error() should return message, got %q", te.Error())
	}
	if te.Unwrap() != nil {
		t.Fatalf("message-only exception has no cause")
	}
}

// TestHttpActionError covers both branches of HttpAction.Error.
func TestHttpActionError(t *testing.T) {
	if (&HttpAction{Code: 302, Message: "redirect"}).Error() != "HttpAction 302: redirect" {
		t.Fatalf("HttpAction with message formatting wrong")
	}
	if (&HttpAction{Code: 401}).Error() != "HttpAction 401" {
		t.Fatalf("HttpAction without message formatting wrong")
	}
}
