package egos

import "errors"

var (
	// ErrConcurrencyViolation is returned when an event store detects a concurrency violation
	ErrConcurrencyViolation = errors.New("concurrency violation")
)
