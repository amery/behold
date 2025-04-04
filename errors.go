package behold

import (
	"errors"

	"darvaza.org/core"
)

// ErrNilReceiver is an error indicating a nil receiver was passed to a method
var ErrNilReceiver = core.ErrNilReceiver

// ErrInvalid is an error indicating an invalid state or argument
var ErrInvalid = core.ErrInvalid

// ErrClosed is an error indicating that an object or resource is closed and cannot be used
var ErrClosed = errors.New("closed")

// ErrReadOnlyTx is an error indicating an attempt to modify a read-only transaction
var ErrReadOnlyTx = errors.New("read-only transaction")
