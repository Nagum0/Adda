package tests

import (
	"adda/pkg"
	"adda/pkg/commands"
	"adda/pkg/db"
	"os"
	"testing"
)

func TestFreshBranch(t *testing.T) {
    // -- SETUP:
    defer os.RemoveAll(".adda")
    defer os.Remove("test.txt")

    commands.Init()

    file, _ := os.Create("test.txt")
    defer file.Close()
    file.WriteString("Hello, world!")

    commands.Add("test.txt")
    commands.Commit("Test commit.")

    commands.Branch("test_branch")

    // -- TESTING:
    _, err := os.Stat(pkg.REFS_HEADS_PATH + "test_branch")
    if os.IsNotExist(err) {
        t.Error(err.Error())
    }

    commitHash, _ := db.ReadRefHead("master")
    newBranchCommitHash, _ := db.ReadRefHead("test_branch")
    if commitHash != newBranchCommitHash {
        t.Errorf("Expected hash: %v; Received hash: %v;", commitHash, newBranchCommitHash)
    }
}
