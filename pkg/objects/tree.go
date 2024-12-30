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
    
    trees := map[string]TreeObject{
        ".": *NewTreeObject(),
    }

    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        fmt.Println(dirs)
        
        // This means that this is a file at the root
        if len(dirs) == 1 {
            root := trees["."]
            root.Files = append(root.Files, dirs[0])
            trees["."] = root
        // This means that we are one subdirectory deep from the root
        } else if len(dirs) == 2 {
            // Appending the subdirectory as the child of the root directory
            root := trees["."]
            root.Children = append(root.Children, dirs[0])
            trees["."] = root
            
            // Appending the file to the subdirectory and create the subdirectory if it doesn't exist
            dir, ok := trees[dirs[0]]
            if !ok {
                dirTree := NewTreeObject()
                dirTree.Files = append(dirTree.Files, dirs[1])
                trees[dirs[0]] = *dirTree
            } else {
                dir.Files = append(dir.Files, dirs[1])
                trees[dirs[0]] = dir
            }
        // The file is deeply nested
        } else {
            currentDir := "."

            // [sub dub file2.txt]
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

    }
    
    for key, val := range trees {
        fmt.Printf("%v:\n  %v\n  %v\n", key, val.Files, val.Children)
    }

    return "nutin", nil
}

func contains[T comparable](xs []T, x T) bool {
    for _, val := range xs {
        if x == val {
            return true
        }
    }

    return false
}
