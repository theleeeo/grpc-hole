package templateparse

import "github.com/TheLeeeo/grpc-hole/fieldselector"

type ParseError interface {
	error
	Location() fieldselector.Selection
}

type parseError struct {
	location fieldselector.Selection
	message  string
	err      error
}

// NewParseError creates a new ParseError using a message
func NewParseError(location fieldselector.Selection, msg string) ParseError {
	return &parseError{
		location: location,
		message:  msg,
	}
}

// ParseErrorWrap wraps an error into a ParseError
func ParseErrorWrap(location fieldselector.Selection, err error) ParseError {
	return &parseError{
		location: location,
		err:      err,
	}
}

func (e *parseError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.message
}

func (e *parseError) Location() fieldselector.Selection {
	return e.location
}
