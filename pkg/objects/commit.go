package objects

import (
	"adda/pkg"
	"adda/pkg/db"
	"adda/pkg/errors"
	"fmt"
	"os"
)

// Commit object.
type CommitObject struct {
    // Commit hash.
    Hash           string
    // The hash for the root tree object.
    RootTree       string
    // The hash for the parent commit.
    ParentCommit   string
    // Author name.
    AuthorName     string
    // Author email.
    AuthorEmail    string
    // Committer name.
    CommitterName  string
    // Committer email.
    CommitterEmail string
    // Commit message.
    Message        string
}

// Sets the hash for the commit.
func (c *CommitObject) GenHash() {
    if c.Hash != "" {
        return
    }
    c.Hash = db.GenSHA1([]byte(c.String()))
}

// Writes the commit to the object database.
func (c CommitObject) DBWrite() error {
    if db.HashExists(c.Hash) {
        return nil
    }

    _, err := os.Stat(pkg.OBJECTS_PATH + c.Hash[2:] + "/")
    if os.IsNotExist(err) {
        if err := os.Mkdir(pkg.OBJECTS_PATH + c.Hash[:2] + "/", os.ModePerm); err != nil {
            return errors.NewCommitError(err.Error())
        }
    }
    
    compressedCommit := db.ZlibCompressString(c.String())
    file, err := os.Create(pkg.OBJECTS_PATH + c.Hash[:2] + "/" + c.Hash[2:])
    if err != nil {
        return errors.NewCommitError(err.Error())
    }
    defer file.Close()
    file.Write(compressedCommit)

    return nil
}

func (c CommitObject) String() string {
    return fmt.Sprintf(
        "tree %v\nparent %v\nauthor %v <%v>\ncommitter %v <%v>\n\n%v",
        c.RootTree, 
        c.ParentCommit,
        c.AuthorName, c.AuthorEmail,
        c.CommitterName, c.CommitterEmail,
        c.Message,
    )
}

// TODO: ParseCommitString
// Takes a string and parses into a CommitObject.
func ParseCommitString(s string) (*CommitObject, error) {
    panic("todo")   
}
