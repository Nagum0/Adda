package objects

import (
	"fmt"
	"strings"
)

type Node struct {
    path     string
    files    []string
    children []*Node
}

func CreateTree() (*Node, error) {
    root := Node{ path: ".", files: make([]string, 0), children: make([]*Node, 0) }

    indexFile, err := ParseIndex()
    if err != nil {
        return nil, err
    }

    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")
        fmt.Println(dirs)

        if len(dirs) == 1 {
            file := dirs[0]
            root.AppendFile(file, ".")
        } else {
            for idx, dir := range dirs {
                if idx == len(dirs) - 1 {
                    fmt.Println("here 1")
                    root.AppendFile(dir, dirs[idx - 1])  
                } else if idx == 2 {
                    fmt.Println("here 2")
                    root.AppendNode(dir, ".")
                } else if idx > 2 {
                    fmt.Println("here 3")
                    root.AppendNode(dir, dirs[idx - 1])
                }
            }
        }
    }

    return &root, nil
}

func (node *Node) AppendNode(dir string, parent string) {
    fmt.Println(dir, parent)

    if node.path == parent && node.ContainsChild(dir) {
        return 
    }
    
    for _, child := range node.children {
        fmt.Println(child.path)
        if child.path == parent && !child.ContainsChild(dir) {
            dirNode := Node {
                path: dir,
                files: make([]string, 0),
                children: make([]*Node, 0),
            }
            child.children = append(child.children, &dirNode)
        } else {
            child.AppendNode(dir, parent)
        }
    }
}

func (node *Node) AppendFile(fileName string, dir string) {
    if node.path == dir {
        node.files = append(node.files, fileName)
    } else {
        for _, child := range node.children {
            child.AppendFile(fileName, dir)
        }
    }
}

func (node Node) ContainsChild(dir string) bool {
    for _, val := range node.children {
        if val.path == dir {
            return true
        }
    }

    return false
}

func (node Node) String() string {
    result := fmt.Sprintf("Path: %v; Children: %v; Files: %v\n", node.path, node.children, node.files)
    
    for _, child := range node.children {
        result += child.String()
    }

    return result
}
