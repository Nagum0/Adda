package objects

import (
	"fmt"
	"strings"
)

// A blob object in the tree object. Holds the hash and the original file name of the blob.
type TreeBlob struct {
    FileName string
    Hash     string
}

func NewTreeBlob(fileName string, hash string) *TreeBlob {
    return &TreeBlob {
        FileName: fileName,
        Hash: hash,
    }
}

// A tree object holds the blob objects of a directory and it's hash.
// Tree object formatting: <object type> <hash> <file name>. The object type can be 0 for a blob
// representing a file or 1 for a subtree representing a subdirectory.
type TreeObject struct {
    // Name of the directory the tree is representing.
    DirName string
    // The hash of the tree object.
    Hash    string
    // The tree's subdirectories
    SubDirs []TreeObject
    // The tree's blob objects   
    Blobs   []TreeBlob
}

func NewTreeObject(dirName string) *TreeObject {
    return &TreeObject {
        DirName: dirName,
        SubDirs: make([]TreeObject, 0),
        Blobs: make([]TreeBlob, 0),
    }
}

// Checks whether the tree object contains the given directory name as one of it's subdirectories.
func (tree TreeObject) ContainsSubDir(dirName string) bool {
    for _, subDir := range tree.SubDirs {
        if dirName == subDir.DirName {
            return true
        }
    }

    return false
}

// Snapshot of the directory structure at the time of a given commit. 
// Maps the directory names to tree objects.
type Snapshot map[string]TreeObject

func NewSnapshot() *map[string]TreeObject {
    return &map[string]TreeObject{}
}

func TakeSnapshot(indexFile Index) Snapshot {
    fmt.Println("INDEX:\n" + indexFile.String())
    
    snapshot := *NewSnapshot()
    snapshot["."] = *NewTreeObject(".")
    
    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        currentDir := "."

        for i := 0; i < len(dirs) - 1; i++ {
            nextDir := dirs[i]

            if _, ok := snapshot[nextDir]; !ok {
                snapshot[nextDir] = *NewTreeObject(nextDir)
            }

            // Add the nextDir to the currentDir as a child if it doesn't already contain it
            if currentDirTree := snapshot[currentDir]; !currentDirTree.ContainsSubDir(nextDir) {
                currentDirTree.SubDirs = append(currentDirTree.SubDirs, *NewTreeObject(nextDir))
                snapshot[currentDir] = currentDirTree
            }

            currentDir = nextDir
        }

        fileName := dirs[len(dirs) - 1]
        currentDirTree := snapshot[currentDir]
        currentDirTree.Blobs = append(currentDirTree.Blobs, *NewTreeBlob(fileName, entry.Hash))
        snapshot[currentDir] = currentDirTree
    }

    for dirName, tree := range snapshot {
        fmt.Printf("%v:\n  %v\n  %v\n", dirName, tree.Blobs, tree.SubDirs)
    }
        
    return snapshot
}

// func generateTreeHash(trees map[string]Snapshot, dirName string, treeHashes *map[string]string) {
//     for _, child := range trees[dirName].Dirs {
//         if _, ok := (*treeHashes)[child]; !ok {
//             generateTreeHash(trees, child, treeHashes)
//         }
//     }
//     
//     (*treeHashes)[dirName] = "<hash>"
// }
