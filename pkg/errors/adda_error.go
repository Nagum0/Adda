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
    return "[INIT ERROR] " + e.msg
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
    return "[ADD ERROR] " + e.msg
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
    return "[BLOB ERROR] " + e.msg
}

// Error representing failure during work with the INDEX file.
type IndexError struct {
    msg string
}

// Creates a new IndexError.
func NewIndexError(msg string) *IndexError {
    return &IndexError{ msg: msg }
}

func (e *IndexError) Error() string {
    return "[INDEX ERROR] " + e.msg
}

// Error representing failure during work with the tree object.
type TreeError struct {
    msg string
}

// Creates a new TreeError.
func NewTreeError(msg string) *TreeError {
    return &TreeError{ msg: msg }
}

func (e *TreeError) Error() string {
    return "[TREE ERROR] " + e.msg
}

// Error representing failure during a commit.
type CommitError struct {
    msg string
}

// Creates a new CommitError.
func NewCommitError(msg string) *CommitError {
    return &CommitError{ msg: msg }
}

func (e *CommitError) Error() string {
    return "[COMMIT ERROR] " + e.msg
}

// Error representing failure during work with branches.
type BranchError struct {
    msg string
}

// Creates a new BranchError.
func NewBranchError(msg string) *BranchError {
    return &BranchError{ msg: msg }
}

func (e *BranchError) Error() string {
    return "[BRANCH ERROR] " + e.msg
}
