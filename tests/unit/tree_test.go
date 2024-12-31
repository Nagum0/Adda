package tests

import (
	"adda/pkg/objects"
	"strings"
	"testing"
)

func TestTakeSnapshot(t *testing.T) {
    indexFile := createTestIndexFile()
    snapshot := objects.TakeSnapshot(indexFile)
    
    // Test the root directory
    rootTree, ok := snapshot["."]

    if !ok {
        t.Error("Root tree not found in snapshot")
    }
    
    if  !rootTree.ContainsSubDir("sub") || 
    !rootTree.ContainsSubDir("stf") ||
    !rootTree.ContainsBlobFile("hello.txt") ||
    !rootTree.ContainsBlobFile("lol.txt") ||
    !rootTree.ContainsBlobFile("lmao.txt") {
        t.Errorf(
            "Root tree does not contain the expected elements\nExpected blobs: hello.txt, lol.txt, lmao.txt; Received blobs: %v\nExpected subdirs: sub, stf; Received subdirs: %v",
            rootTree.Blobs,
            rootTree.SubDirs,
        )
    }

    if rootTree.Hash != "bc8c7392523d88fdc04de7b48168bc51864ac8b8" {
        t.Errorf("Expected root tree hash: bc8c7392523d88fdc04de7b48168bc51864ac8b8; Received root tree hash: %v", rootTree.Hash)
    }

    // Test the sub directory
    subTree, ok := snapshot["sub"]

    if !ok {
        t.Error("Sub tree not found in snapshot")
    }

    if  !subTree.ContainsSubDir("dub") || 
    !subTree.ContainsBlobFile("file.txt") {
        t.Errorf(
            "Root tree does not contain the expected elements\nExpected blobs: [file.txt]; Received blobs: %v\nExpected subdirs: [dub]; Received subdirs: %v",
            subTree.Blobs,
            subTree.SubDirs,
        )
    }

    if subTree.Hash != "0483d54712d4c601af4c5d36ca603417a33d8b23" {
        t.Errorf("Expected sub tree hash: 0483d54712d4c601af4c5d36ca603417a33d8b23; Received sub tree hash: %v", subTree.Hash)
    }

    // Test the dub directory
    dubTree, ok := snapshot["dub"]

    if !ok {
        t.Error("Sub tree not found in snapshot")
    }

    if !dubTree.ContainsBlobFile("file2.txt") {
        t.Errorf(
            "Root tree does not contain the expected elements\nExpected blobs: [file2.txt]; Received blobs: %v\nExpected dubdirs: []; Received dubdirs: %v",
            dubTree.Blobs,
            dubTree.SubDirs,
        )
    }

    if dubTree.Hash != "a2e863286e13f3a9845a13f8a6b1c87604235cf7" {
        t.Errorf("Expected dub tree hash: a2e863286e13f3a9845a13f8a6b1c87604235cf7; Received dub tree hash: %v", dubTree.Hash)
    }

    // Test the stf directory
    stfTree, ok := snapshot["stf"]

    if !ok {
        t.Error("Stf tree not found in snapshot")
    }

    if !stfTree.ContainsBlobFile("s.txt") {
        t.Errorf(
            "Root tree does not contain the expected elements\nExpected blobs: [s.txt]; Received blobs: %v\nExpected stfdirs: []; Received stfdirs: %v",
            stfTree.Blobs,
            stfTree.SubDirs,
        )
    }

    if stfTree.Hash != "f615d1efcb4ce91b1cacab5d77dd5890f31b7537" {
        t.Errorf("Expected stf tree hash: f615d1efcb4ce91b1cacab5d77dd5890f31b7537; Received stf tree hash: %v", stfTree.Hash)
    }
}

func createTestIndexFile() objects.Index {
    testIndexFileString := `0  7489b43a3a5afe97c4d388295fd65051d4cc235d     lol.txt
    0  9c1ef3eac3179bb546eb9f270052caadb5ba7cf3     sub/file.txt
    0  5ee36c04ffae7fb264e2ac08af413e780499d810     sub/dub/file2.txt
    0  2ef050e967ad03bfdcf220400a4b2ac1caa9fc57     lmao.txt
    0  0bbea09c392b5b9c8e5d113d55c93e0e5143163c     stf/s.txt
    0  36b2dbe1d6b45b568b6f772419ef24ca8fe2866f     hello.txt`

    testIndexFile := objects.Index {
        Entries: make(map[string]objects.Entry),
    }

    for _, entryString := range strings.Split(testIndexFileString, "\n") {
        splitEntryStr := strings.Fields(entryString)
        entry := objects.Entry{
            FilePath: splitEntryStr[2],
            Hash: splitEntryStr[1],
        }

        switch splitEntryStr[0] {
            case "0":
                entry.FileType = objects.FILE
                break
            case "1":
                entry.FileType = objects.DIR
                break
            default:
                break
        }

        testIndexFile.Entries[splitEntryStr[2]] = entry
    }

    return testIndexFile
}