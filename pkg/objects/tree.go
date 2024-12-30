package objects

import (
	"fmt"
	"strings"
)

// Tree objects are snapshots of the directory structure at the time of a given commit.
// Tree object formatting: <object type> <hash> <file name>. The object type can be 0 for a blob
// representing a file or 1 for a subtree representing a subdirectory.
type TreeObject struct {
    Files    []string
    Children []string
}

func NewTreeObject() *TreeObject {
    return &TreeObject {
        Files: make([]string, 0),
        Children: make([]string, 0),
    }
}

func IndexToTree(indexFile Index) (string, error) {
    fmt.Println("INDEX:\n" + indexFile.String())
    
    // Map the directory names to tree objects
    trees := map[string]TreeObject{
        ".": *NewTreeObject(),
    }

    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        currentDir := "."

        for i := 0; i < len(dirs) - 1; i++ {
            nextDir := dirs[i]

            if _, ok := trees[nextDir]; !ok {
                trees[nextDir] = *NewTreeObject()   
            }

            // Add the nextDir to the currentDir as a child if it doesn't already contain it
            if currentDirTree := trees[currentDir]; !contains(currentDirTree.Children, nextDir) {
                currentDirTree.Children = append(currentDirTree.Children, nextDir)
                trees[currentDir] = currentDirTree
            }

            currentDir = nextDir
        }

        fileName := dirs[len(dirs) - 1]
        currentDirTree := trees[currentDir]
        currentDirTree.Files = append(currentDirTree.Files, fileName)
        trees[currentDir] = currentDirTree
    }

    for key, val := range trees {
        fmt.Printf("%v:\n  %v\n  %v\n", key, val.Files, val.Children)
    }

    // Map the directory names to tree hashes
    treeHashes := map[string]string{}
    generateTreeHash(trees, ".", &treeHashes)
    
    for k, v := range treeHashes {
        fmt.Printf("%v -> %v\n", k, v)
    }

    // When we have the hashes and the contents we make a snapshot thingy !!!
    
    return "nutin", nil
}

func generateTreeHash(trees map[string]TreeObject, dirName string, treeHashes *map[string]string) {
    for _, child := range trees[dirName].Children {
        if _, ok := (*treeHashes)[child]; !ok {
            generateTreeHash(trees, child, treeHashes)
        }
    }
    
    (*treeHashes)[dirName] = "<hash>"
}

func contains[T comparable](xs []T, x T) bool {
    for _, val := range xs {
        if x == val {
            return true
        }
    }

    return false
}
