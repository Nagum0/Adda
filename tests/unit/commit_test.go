package tests

import (
	"adda/pkg/db"
	"adda/pkg/objects"
	"strings"
	"testing"
)

func TestParseCommitString(t *testing.T) {
    // -- SETUP:
    var commitString string = `
tree a6f410699fdf6949c49ab39f6bd29e4dacb3508d
parent caf9f837b68cfd5e728cc9a5cfeac8f5b3704cd5
author tester test@email.com
committer tester test@email.com

Second
    `
    trimmedString := strings.Trim(commitString, "\n")
    commitObjectHash := db.GenSHA1([]byte(trimmedString))
    commitObject, err := objects.ParseCommitString(trimmedString)
    
    // -- TEST:
    if err != nil {
        t.Error(err.Error())
    }
    
    if commitObject.Hash != commitObjectHash {
        t.Errorf("Expected commit object hash: %v; Received commit object hash: %v", commitObjectHash, commitObject.Hash)
    }

    if commitObject.ParentCommit != "caf9f837b68cfd5e728cc9a5cfeac8f5b3704cd5" {
        t.Errorf("Expected parent commit hash: caf9f837b68cfd5e728cc9a5cfeac8f5b3704cd5; Received parent commit hash: %v", commitObject.ParentCommit)
    }

    if commitObject.RootTree != "a6f410699fdf6949c49ab39f6bd29e4dacb3508d" {
        t.Errorf("Expected root tree hash: a6f410699fdf6949c49ab39f6bd29e4dacb3508d; Received root tree hash: %v", commitObject.RootTree)
    }

    if commitObject.AuthorName != "tester" || commitObject.AuthorEmail != "test@email.com" {
        t.Errorf("Expected author info: tester test@email.com; Received author info: %v %v", commitObject.AuthorName, commitObject.AuthorEmail)
    }

    if commitObject.CommitterName != "tester" || commitObject.CommitterEmail != "test@email.com" {
        t.Errorf("Expected committer info: tester test@email.com; Received committer info: %v %v", commitObject.CommitterName, commitObject.CommitterEmail)
    }
    
    if commitObject.Message != "Second" {
        t.Errorf("Expected commit message: Second; Received commit message: %v", commitObject.Message)
    }
}
