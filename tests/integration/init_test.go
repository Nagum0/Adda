package tests

import (
	"adda/pkg/commands"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
    defer os.RemoveAll(".adda/")

    if err := commands.Init(); err != nil {
        t.Error(err.Error())
    }

    initFiles := []string{".adda/", ".adda/objects", ".adda/branches", ".adda/refs", ".adda/refs/heads", ".adda/INDEX", ".adda/HEAD"}
    
    for _, initFile := range initFiles {
        if _, err := os.Stat(initFile); err != nil && os.IsNotExist(err) {
            t.Errorf("%v was not created successfully", initFile)
        }
    }
}
