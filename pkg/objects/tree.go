package objects

import (
	"adda/pkg/errors"
	"fmt"
	"strings"
)

// A node inside the tree object structure.
type Node struct {
    // The name of the directory the node represents.
    Path     string
    // The files in the directory.
    Files    []string
    // The subdirectories represented as a slice of pointers to Nodes.
    Children []*Node
}

// Initializes a new node.
func NewNode(Path string) *Node {
    return &Node {
        Path: Path,
        Files: make([]string, 0),
        Children: make([]*Node, 0),
    }
}

// Parses the INDEX file into a tree. If it fails it will return an IndexError.
func IndexToTree() (*Node, error) {
    root := Node{ Path: ".", Files: make([]string, 0), Children: make([]*Node, 0) }

    indexFile, err := ParseIndex()
    if err != nil {
        return nil, err
    }

    for _, entry := range indexFile.Entries {
        dirs := strings.Split(entry.FilePath, "/")

        if len(dirs) == 1 {
            file := dirs[0]
            root.AppendFile(file, ".")
        } else {
            for idx, dir := range dirs {
                if idx == 0 {
                    root.AppendNode(dir, ".")
                } else if idx == len(dirs) - 1 {
                    root.AppendFile(dir, dirs[idx - 1])
                } else if idx > 0 {
                    root.AppendNode(dir, dirs[idx - 1])
                }
            }
        }
    }

    return &root, nil
}

// Appends a new node into the tree object.
// The node represents a directory. 
// The directory name and the parent directory name must be given as paremeters.
func (node *Node) AppendNode(dir string, parent string) {
    if node.Path == parent && !node.ContainsChild(dir) {
        dirNode := Node {
            Path: dir,
            Files: make([]string, 0),
            Children: make([]*Node, 0),
        }
        node.Children = append(node.Children, &dirNode)
    } else {
        for _, child := range node.Children {
            child.AppendNode(dir, parent)
        }
    }
}

// Appends a file to the given directory.
func (node *Node) AppendFile(fileName string, dir string) {
    if node.Path == dir {
        node.Files = append(node.Files, fileName)
    } else {
        for _, child := range node.Children {
            child.AppendFile(fileName, dir)
        }
    }
}

// Checks whether the node the function was called on contains the given directory name
// in it's children slice.
func (node Node) ContainsChild(dir string) bool {
    for _, val := range node.Children {
        if val.Path == dir {
            return true
        }
    }

    return false
}

// Searches for the given directory inside the tree and returns it as a *Node.
func (node Node) GetNode(dir string) (*Node, error) {
    if node.Path == dir {
        return &node, nil
    }

    for _, child := range node.Children {
        return child.GetNode(dir)
    }

    return nil, errors.NewTreeError(fmt.Sprintf("Node: %v not found.", dir))
}

func (node Node) String() string {
    result := fmt.Sprintf("Path: %v; Files: %v; Children: %v", node.Path, node.Files, node.Children)
    return result
}
