package git

import "fmt"

// GitError wraps a git command failure with its output.
type GitError struct {
	cmd    string
	output string
	err    error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("%s failed: %s\n%s", e.cmd, e.err, e.output)
}

func (e *GitError) Unwrap() error { return e.err }
