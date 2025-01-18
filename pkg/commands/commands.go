package commands

import (
	"adda/pkg/db"
	"adda/pkg/errors"
	"adda/pkg/objects"
	"fmt"
	"os"
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

// Creates a commit object and writes it to the object database. Also sets the 
// refs/heads/<current branch> to the new commit object's hash.
func Commit(msg string) error {
    indexFile, err := objects.ParseIndex()
    if err != nil {
        return errors.NewCommitError(err.Error())
    }
    
    snapshot := objects.TakeSnapshot(*indexFile)
    if err = snapshot.WriteSnapshotToDatabase(); err != nil {
        return errors.NewCommitError(err.Error())
    }
        
    return nil
}
