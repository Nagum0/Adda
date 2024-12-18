package commands

import (
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

func Add(filePath string) error {
    blob := objects.NewBlob(filePath, objects.FILE)
    _, err := blob.GenerateContents()
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(blob.Contents)
    return nil
}
