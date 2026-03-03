package commands

import "errors"

// ErrAbort is returned when the user cancels a confirmation prompt.
// The main function maps this to exit code 1.
var ErrAbort = errors.New("aborted")

// errAbort is the unexported alias used within commands.
var errAbort = ErrAbort
