package tests

import (
	"adda/pkg/commands"
	"adda/pkg/index"
	"adda/pkg/objects"
	"fmt"
	"os"
	"testing"
)

func TestAddNewFile(t *testing.T) {
    defer os.RemoveAll(".adda")
    defer os.Remove("test.txt")

    if err := commands.Init(); err != nil {
        t.Error(err.Error())
    }

    file, err := os.Create("test.txt")
    if err != nil {
        t.Error(err.Error())
    }
    file.WriteString("Hello, World!")
    file.Close()
    
    blob := objects.NewBlob("test.txt", objects.FILE)
    hash, _ := blob.Hash()

    if err = commands.Add("test.txt"); err != nil {
        t.Error(err.Error())
    }
    
    indexFile, err := index.ParseIndex()
    if err != nil {
        t.Error(err.Error())
    }

    // Check if the INDEX file has the entry
    entry := indexFile.Entries["test.txt"]
    if entry.String() != fmt.Sprintf("0  %v\ttest.txt", hash) {
        t.Errorf("Entry for added file not found in INDEX: %v -> %v", hash, entry.Hash)
    }

    // Check if the blob object was created
    if _, err = os.Stat(".adda/objects/" + hash[:2] + "/" + hash[2:]); err != nil && os.IsNotExist(err) {
        t.Errorf("Blob file for hash: %v/%v not found", hash[0:2], hash[2:])
    }
}

func TestAddChangedFile(t *testing.T) {

}
