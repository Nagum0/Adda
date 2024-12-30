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

func (tree TreeObject) ContainsSubDir(dirName string) bool {
    for _, subDir := range tree.SubDirs {
        if dirName == subDir.DirName {
            return true
        }
    }

    return false
}

// Snapshot of the directory structure at the time of a given commit. 
type Snapshot struct {
    // The directories of the snapshot.
    Dirs map[string]TreeObject
}

func NewSnapshot() *Snapshot {
    return &Snapshot {
        Dirs: make(map[string]TreeObject, 0),
    }
}

func IndexToTree(indexFile Index) (string, error) {
    fmt.Println("INDEX:\n" + indexFile.String())
    
    snapshot := *NewSnapshot()
    snapshot.Dirs["."] = *NewTreeObject(".")
    
    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        currentDir := "."

        for i := 0; i < len(dirs) - 1; i++ {
            nextDir := dirs[i]

            if _, ok := snapshot.Dirs[nextDir]; !ok {
                snapshot.Dirs[nextDir] = *NewTreeObject(nextDir)
            }

            // Add the nextDir to the currentDir as a child if it doesn't already contain it
            if currentDirTree := snapshot.Dirs[currentDir]; !currentDirTree.ContainsSubDir(nextDir) {
                currentDirTree.SubDirs = append(currentDirTree.SubDirs, *NewTreeObject(nextDir))
                snapshot.Dirs[currentDir] = currentDirTree
            }

            currentDir = nextDir
        }

        fileName := dirs[len(dirs) - 1]
        currentDirTree := snapshot.Dirs[currentDir]
        currentDirTree.Blobs = append(currentDirTree.Blobs, *NewTreeBlob(fileName, entry.Hash))
        snapshot.Dirs[currentDir] = currentDirTree
    }

    for dirName, tree := range snapshot.Dirs {
        fmt.Printf("%v:\n  %v\n  %v\n", dirName, tree.Blobs, tree.SubDirs)
    }

    // // Map the directory names to tree hashes
    // treeHashes := map[string]string{}
    // generateTreeHash(snapshot.Dirs, ".", &treeHashes)
    // 
    // for k, v := range treeHashes {
    //     fmt.Printf("%v -> %v\n", k, v)
    // }

    // When we have the hashes and the contents we make a snapshot thingy !!!
    
    return "nutin", nil
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

func contains[T comparable](xs []T, x T) bool {
    for _, val := range xs {
        if x == val {
            return true
        }
    }

    return false
}
