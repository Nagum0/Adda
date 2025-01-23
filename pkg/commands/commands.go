package commands

import (
	"adda/pkg/db"
	"adda/pkg/errors"
	"adda/pkg/objects"
	"fmt"
	"os"
	"strings"
)

// Initializes an Adda repository by creating all of the neccessary init files.
// Return an InitError if the process failes.
func Init() error {
    // Creating .adda/ directory
    err := os.Mkdir(".adda", os.ModePerm)
    if err != nil {
        return errors.NewInitError(fmt.Sprintf("Error while initializing adda repository: %v", err.Error()))
    }
    
    // Creating the needed subdirectories
    dirs := []string{"objects", "refs", "branches", "refs/heads"}
    for _, dir := range dirs {
        err = os.MkdirAll(".adda/" + dir, os.ModePerm)
        if err != nil {
            return errors.NewInitError(fmt.Sprintf("Error while creating %v directory: %v", dir, err.Error()))
        }
    }

    // Creating the needed files in root
    head, err := os.Create(".adda/HEAD")
    if err != nil {
        return errors.NewInitError("Error while creating HEAD file")
    }
    defer head.Close()

    index, err := os.Create(".adda/INDEX")
    if err != nil {
        return errors.NewInitError("Error while creating INDEX file")
    }
    defer index.Close()

    return nil
}

// Generate a blob object and write it to the filesystem.
// The blob object will be created at .adda/objects/<hash_begin>/<hash_rest>.
// This will also update the INDEX and map the hash of the file to the filepath.
func Add(filePath string) error {
    blob := objects.NewBlob(filePath, objects.FILE)
    hash, err := blob.Hash()
    if err != nil {
        return errors.NewAddError(err.Error(), filePath)
    }

    if db.HashExists(hash) {
        return nil
    }

    err = blob.WriteBlob()
    if err != nil {
        return errors.NewAddError(err.Error(), filePath)
    }

    indexFile, err := objects.ParseIndex()
    if err != nil {
        return errors.NewAddError(err.Error(), filePath)
    }

    err = indexFile.Update(*blob)
    if err != nil {
        return errors.NewAddError(err.Error(), filePath)
    }

    return nil
}

// TODO: ONLY COMMIT IF SOMETHING HAS ACTUALY BEEN MODIFIED!
// Creates a commit object and writes it to the object database. Also sets the 
// refs/heads/<current branch> to the new commit object's hash.
func Commit(msg string) error {
    indexFile, err := objects.ParseIndex()
    if err != nil {
        return errors.NewCommitError(err.Error())
    }
    
    // Take snapshot of the current state and save it in the database
    snapshot := objects.TakeSnapshot(*indexFile)
    if err = snapshot.DBWrite(); err != nil {
        return errors.NewCommitError(err.Error())
    }
    
    commit := objects.CommitObject{
        Message: msg,
        RootTree: snapshot["."].Hash,
        AuthorName: "tester",
        AuthorEmail: "test@email.com",
        CommitterName: "tester",
        CommitterEmail: "test@email.com",
    }

    head, err := db.ReadHEAD()
    if err != nil {
        return errors.NewCommitError(err.Error())
    }
    
    // Get the current branch
    var currentBranch string

    if head == "" {
        if err := db.SetHEAD("refs/heads/master"); err != nil {
            return errors.NewCommitError(err.Error())
        }
        currentBranch = "master"
    } else {
        currentBranch = strings.Split(head, "/")[2]
        currentBranch = strings.Trim(currentBranch, "\n")
    }

    // Get the parent hash if it exists
    refHead, err := db.ReadRefHead(currentBranch)
    if err != nil {
        if os.IsNotExist(err) {
            // TODO: If the ref head is empty it must mean that we're on fresh branch
            // and we should set the parent commit to be the last commit we branched off of.
            commit.ParentCommit = ""
        } else {
            return errors.NewCommitError(err.Error())
        }
    }
    commit.ParentCommit = strings.Trim(refHead, "\n")
    
    // Generate the commit hash
    commit.GenHash()

    // Write the commit object to the object database
    if err := commit.DBWrite(); err != nil {
        return err
    }
    
    // Set the refs/heads/currentBranch to the commit's hash
    if err := db.SetRefsHead(currentBranch, commit.Hash); err != nil {
        return errors.NewCommitError(err.Error())
    }

    fmt.Println(commit.String())
        
    return nil
}

// Prints the contents of the file with the given hash.
func Cat(hash string) error {
    data, err := db.DBRead(hash)
    if err != nil {
        return errors.NewBlobError("[CAT ERROR]")
    }
    
    decompressedData, err := db.Decompress(data)
    if err != nil {
        return errors.NewBlobError("[CAT ERROR]")
    }

    fmt.Println(string(decompressedData))

    return nil
}
