package objects

import (
	"adda/pkg"
	"adda/pkg/db"
	"adda/pkg/errors"
	"fmt"
	"os"
	"strings"
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
        "tree %v\nparent %v\nauthor %v %v\ncommitter %v %v\n\n%v",
        c.RootTree, 
        c.ParentCommit,
        c.AuthorName, c.AuthorEmail,
        c.CommitterName, c.CommitterEmail,
        c.Message,
    )
}

// Takes a string and parses it into a CommitObject.
func ParseCommitString(s string) (*CommitObject, error) {
    commitObject := CommitObject{}
    commitObject.Hash = db.GenSHA1([]byte(s))
    s = strings.Replace(s, "\n\n", "\n", -1)
    lines := strings.Split(s, "\n")
    lines = lines[:len(lines) - 1]

    for _, line := range lines {
        lineParts := strings.Fields(line)

        if len(lineParts) == 1 {
            commitObject.Message = line
        } else if len(lineParts) == 2 {
            switch lineParts[0] {
            case "tree":
                commitObject.RootTree = lineParts[1]
                break
            case "parent":
                commitObject.ParentCommit = lineParts[1]
                break
            default:
                return nil, errors.NewCommitError("Illegal commit format.")
            }
        } else if len(lineParts) == 3 {
            switch lineParts[0] {
            case "author":
                commitObject.AuthorName = lineParts[1]
                commitObject.AuthorEmail = lineParts[2]
                break
            case "committer":
                commitObject.CommitterName = lineParts[1]
                commitObject.CommitterEmail = lineParts[2]
                break
            default:
                return nil, errors.NewCommitError("Illegal commit format.")
            }
        } else {
            fmt.Println("here")
            return nil, errors.NewCommitError("Illegal commit format.")
        }
    }

    return &commitObject, nil
}
