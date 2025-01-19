package tests

import (
	"adda/pkg/commands"
	"adda/pkg/db"
	"os"
	"testing"
)

func TestCommitFresh(t *testing.T) {
    // ---------- SETUP ----------
    defer os.RemoveAll(".adda/")
    defer os.Remove("test.txt")

    if err := commands.Init(); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    f, err := os.Create("test.txt")
    if err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }
    defer f.Close()
    if _, err := f.WriteString(Lorem1); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    if err := commands.Add("test.txt"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    if err := commands.Commit("Init"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    // ---------- CHECKS ----------
    head, _ := db.ReadHEAD()
    if head != "refs/heads/master" {
        t.Errorf("Expected: refs/heads/master; Received: %v", head)
    }
    
    expectedMasterHead := "90c3aa7546289675025508a29b84ff926af4859e"
    masterHead, _ := db.ReadRefHead("master")
    if masterHead != expectedMasterHead {
        t.Errorf("Expected: %v; Received: %v", expectedMasterHead, masterHead)
    }

    if !db.HashExists(masterHead) {
        t.Errorf("Commit object could not be found at hash: %v", masterHead)
    }
}

func TestNewCommit(t *testing.T) {
    // ---------- SETUP ----------
    defer os.RemoveAll(".adda/")
    defer os.Remove("test.txt")

    if err := commands.Init(); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    f, err := os.Create("test.txt")
    if err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }
    defer f.Close()
    if _, err := f.WriteString(Lorem1); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    if err := commands.Add("test.txt"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }
    
    // First commit
    if err := commands.Commit("Init"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    // Change test file
    if _, err := f.WriteString(Lorem2); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }
    
    // Second add
    if err := commands.Add("test.txt"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }

    // Second commit
    if err := commands.Commit("Second"); err != nil {
        t.Errorf("Error while setting up test: %v", err.Error())
    }
    
    // ---------- CHECKS ----------
    expectedRefHead := "dbd0e459b8a446918cec28f6e2b3f5385dcad11b"
    refHead, _ := db.ReadRefHead("master")

    if expectedRefHead != refHead {
        t.Errorf("Expected: %v; Received: %v", expectedRefHead, refHead)
    }

    if !db.HashExists(refHead) {
        t.Errorf("Could not find commit object at hash: %v", refHead)
    }
}

// TODO: This test should be written after branching is implemented.
func TestCommitFromDifferentBranch(t *testing.T) {

}
