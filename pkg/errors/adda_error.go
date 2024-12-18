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

// Error representing a failure during adding a file to the stagin area (INDEX).
type AddError struct {
    msg      string
    filePath string
}

// Creates a new AddError.
func NewAddError(msg string, filePath string) *AddError {
    return &AddError{ msg: msg, filePath: filePath }
}

func (e *AddError) Error() string {
    return e.msg
}

// Error representing a failure during blob object creation.
type BlobError struct {
    msg string
}

// Creates a new BlobError.
func NewBlobError(msg string) *BlobError {
    return &BlobError{ msg: msg }
}

func (e *BlobError) Error() string {
    return e.msg
}
