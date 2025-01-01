package env

import "errors"

// ErrInvalidLogFormat is returned when the provided text cannot be
// unmarshaled to a valid LogFormat.
var ErrInvalidLogFormat = errors.New("invalid log format requested")

// ErrMarshalFailed is returned when an invalid installType is marshaled
// to text.
var ErrMarshalFailed = errors.New("failed to marshal install type")

// ErrUnmarshalFailed is returned when text can't be unmarshaled into
// one of the known installType values.
var ErrUnmarshalFailed = errors.New("failed to unmarshal install type")
