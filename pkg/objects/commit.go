package objects

import "fmt"

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
