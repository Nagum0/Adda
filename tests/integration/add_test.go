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
    files := map[string]string{
        "test1.txt": "Hello, World!",
        "test2.txt": Lorem1,
        "test3.txt": Lorem2,
    }
    hashes := make(map[string]string)

    defer os.RemoveAll(".adda")

    if err := commands.Init(); err != nil {
        t.Error(err.Error())
    }

    for filePath, text := range files {
        defer os.Remove(filePath)

        file, err := os.Create(filePath)
        if err != nil {
            t.Error(err.Error())
        }
        file.WriteString(text)
        file.Close()

        blob := objects.NewBlob(filePath, objects.FILE)
        hash, _ := blob.Hash()
        hashes[filePath] = hash

        if err = commands.Add(filePath); err != nil {
            t.Error(err.Error())
        }

        // Check if the blob object was created
        if _, err = os.Stat(".adda/objects/" + hash[:2] + "/" + hash[2:]); err != nil && os.IsNotExist(err) {
            t.Errorf("Blob file for hash: %v/%v not found", hash[0:2], hash[2:])
        }
    }

    indexFile, err := index.ParseIndex()
    if err != nil {
        t.Error(err.Error())
    }

    // Check if the INDEX file has the entries
    for filePath, hash := range hashes {
        entry := indexFile.Entries[filePath]
        expectedEntry := fmt.Sprintf("0  %v\t%v", hash, filePath)
        
        if entry.String() != expectedEntry {
            t.Errorf("Expected entry: %v. Received entry: %v", expectedEntry, entry.String())
        }
    }
}

func TestAddChangedFile(t *testing.T) {

}
