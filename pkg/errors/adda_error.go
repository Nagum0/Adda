package errors

type InitError struct {
    msg string
}

func NewInitError(msg string) *InitError {
    return &InitError{ msg: msg }
}

func (e *InitError) Error() string {
    return e.msg
}
