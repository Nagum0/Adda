package objects

import (
	"adda/pkg/errors"
	"fmt"
	"os"
)

// Commit object.
type CommitObject struct {
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

func CreateCommitObject(snapshot Snapshot) (*CommitObject, error) {
    head, err := GetHEAD()
    if err != nil {
        return nil, errors.NewCommitError(err.Error())
    }
    
    // Currently branches are not handled so it only creates a master file at .adda/refs/heads/.
    if head == "" {
        if err := initHEAD(); err != nil {
            return nil, err
        }
    }

    commit := CommitObject{
        RootTree: snapshot["."].Hash,
    }

    return &commit, nil
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

func initHEAD() error {
    fmt.Println("HEAD is empty. Creating master branch and setting refs head.")

    // Create the reference file for the master branch
    if err := CreateReferenceFile("master"); err != nil {
        return errors.NewCommitError(err.Error())
    }

    file, err := os.Open(".adda/HEAD")
    if err != nil {
        return errors.NewCommitError(err.Error())
    }
    defer file.Close()
    file.WriteString("refs/heads/master")
    
    return nil
}
