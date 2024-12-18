package errors

// Error representing a failure during repository initialization.
type InitError struct {
    msg string
}

// Creates a new InitError.
func NewInitError(msg string) *InitError {
    return &InitError{ msg: msg }
}

func (e *InitError) Error() string {
    return e.msg
}
