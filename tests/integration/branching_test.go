package tests

import (
	"adda/pkg"
	"adda/pkg/commands"
	"adda/pkg/db"
	"adda/pkg/objects"
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

    // -- TEST:
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

func TestGotoBranchHeadUpdate(t *testing.T) {
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
    // TODO: commands.Goto("test_branch")

    // -- TEST:
    head, _ := db.ReadHEAD()
    if head != "refs/heads/test_branch" {
        t.Errorf("Expected HEAD: refs/heads/test_branch; Received HEAD: %v", head)
    }   
}

func TestParseTreeObjectString(t *testing.T) {
    // -- SETUP:
    defer os.RemoveAll(".adda/")
    defer os.RemoveAll("sub1/")
    defer os.RemoveAll("sub2/")
    defer os.Remove("test.txt")
    test_txt_data := `Line 1: Astolfo
Line 2: Helldivers
Line 3: Ferris
`   
    os.Mkdir("sub1", os.ModePerm)
    os.Mkdir("sub2", os.ModePerm)

    test_txt, _ := os.Create("test.txt")
    defer test_txt.Close()
    test1_txt, _ := os.Create("sub1/test1.txt")
    defer test1_txt.Close()
    test2_txt, _ := os.Create("sub2/test2.txt")
    defer test2_txt.Close()

    test_txt.WriteString(test_txt_data)
    test1_txt.WriteString("Hello, from sub1!\n")
    test2_txt.WriteString("Hello, from sub2!\n")
    
    commands.Init()
    commands.Add("test.txt")
    commands.Add("sub1/test1.txt")
    commands.Add("sub2/test2.txt")
    commands.Commit("Init")

    var rootTreeHash string = "b897d90a452c451227eb193248bf7dc4c014b089"
    var sub1TreeHash string = "295d907a38e794ef47423cbe0d2d9e2939489275"
    var sub2TreeHash string = "935f7c9f96e21e1e8a17c5e1681f2291ebb61ba7"
    var testTxtHash string = "7f2a945e69db252edaabc5cffa8cb4e70494c501"
    var test1TxtHash string = "6870bb2cbd7dd6ab9d75aa5186e6743c4e35d950"
    var test2TxtHash string = "2183046bdbb381739b3c93a3ed4fed37375a9498"

    rootTreeData, _ := db.DBRead(rootTreeHash)
    rootTreeString, _ := db.DecompressToString(rootTreeData)
    rootTreeObject, err := objects.ParseTreeObjectString(rootTreeString, true)

    // -- TEST:
    if err != nil {
        t.Error(err.Error())
    }

    if rootTreeObject.DirName != "." {
        t.Errorf("Expected root tree DirName: .; Received root tree DirName: %v", rootTreeObject.DirName)
    }

    if rootTreeObject.Hash != rootTreeHash {
        t.Errorf("Expected root tree hash: %v; Received root tree hash: %v", rootTreeHash, rootTreeObject.Hash)
    }

    if len(rootTreeObject.Blobs) != 1 {
        t.Errorf("Expected root tree Blobs len: 1; Received root tree Blobs len: %v", len(rootTreeObject.Blobs))
    }

    rootTreeBlob := rootTreeObject.Blobs[0]
    if rootTreeBlob.Hash != testTxtHash {
        t.Errorf("Expected test.txt hash: %v; Received test.txt hash: %v", testTxtHash, rootTreeBlob.Hash)
    }

    if rootTreeBlob.FileName != "test.txt" {
        t.Errorf("Expected: test.txt; Received: %v", rootTreeBlob.FileName)
    }

    if len(rootTreeObject.SubDirs) != 2 {
        t.Errorf("Expected root tree SubDirs len: 2; Received root tree SubDirs len: %v", len(rootTreeObject.SubDirs))
    }
    
    sub1Tree := rootTreeObject.SubDirs[0]
    sub2Tree := rootTreeObject.SubDirs[1]

    if sub1Tree.Hash != sub1TreeHash {
        t.Errorf("Expected sub1 tree hash: %v; Received sub1 tree hash: %v", sub1TreeHash, sub1Tree.Hash)
    }

    if len(sub1Tree.Blobs) != 1 {
        t.Errorf("Expected sub1 tree Blobs len: 1; Received sub1 tree Blobs len: %v", len(sub1Tree.Blobs))
    }
    
    sub1TreeBlob := sub1Tree.Blobs[0]
    if sub1TreeBlob.Hash != test1TxtHash {
        t.Errorf("Expected test1.txt hash: %v; Received test1.txt hash: %v", test1TxtHash, sub1TreeBlob.Hash)
    }

    if sub1TreeBlob.FileName != "test1.txt" {
        t.Errorf("Expected: test1.txt; Received: %v", sub1TreeBlob.FileName)
    }

    if sub2Tree.Hash != sub2TreeHash {
        t.Errorf("Expected sub2 tree hash: %v; Received sub2 tree hash: %v", sub2TreeHash, sub2Tree.Hash)
    }

    if len(sub2Tree.Blobs) != 1 {
        t.Errorf("Expected sub2 tree Blobs len: 1; Received sub2 tree Blobs len: %v", len(sub2Tree.Blobs))
    }

    sub2TreeBlob := sub2Tree.Blobs[0]
    if sub2TreeBlob.Hash != test2TxtHash {
        t.Errorf("Expected test2.txt hash: %v; Received test2.txt hash: %v", test2TxtHash, sub2TreeBlob.Hash)
    }

    if sub2TreeBlob.FileName != "test2.txt" {
        t.Errorf("Expected: test2.txt; Received: %v", sub2TreeBlob.FileName)
    }
}

func TestBuildBranchsLatestSnapshot(t *testing.T) {
    // -- SETUP:

    // -- TEST:
}
