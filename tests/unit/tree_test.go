package tests

import (
	"adda/pkg/objects"
	"fmt"
	"testing"
)

func TestGetNode(t *testing.T) {
    root := objects.NewNode(".")
    root.AppendNode("sub", ".")
    root.AppendNode("dub", ".")
    root.AppendNode("gus", "sub")
    root.AppendNode("foo", "dub")
    root.AppendNode("bar", "gus")

    node, err := root.GetNode("gus")
    if err != nil {
        t.Error(err.Error())
    }

    if node.Path != "gus" {
        t.Errorf("Expected path: gus; Received path: %v", node.Path)
    }
    
    if len(node.Files) != 0 {
        t.Errorf("Expected files length: 0; Received files length: %v", len(node.Files))
    }
    
    barNode := objects.NewNode("bar")
    var child objects.Node
    if len(node.Children) == 1 {
        child = *node.Children[0]
    } else {
        t.Errorf("Expetced children lentgh: 1; Received children length: %v", len(node.Children))
    }

    if child.Path != "bar" || len(child.Files) != 0 || len(child.Children) != 0 {
        t.Errorf("Expected: %v; Received: %v", barNode, child)
    }
}

func TestAppendNode(t *testing.T) {
    root := objects.NewNode(".")
    root.AppendNode("sub", ".")
    root.AppendNode("dub", ".")
    root.AppendNode("gus", "sub")
    root.AppendNode("foo", "dub")
    root.AppendNode("bar", "gus")

    for _, child := range root.Children {
        fmt.Println(child)
    }
}

func TestAppendFile(t *testing.T) {

}

func TestContainsChild(t *testing.T) {

}

func TestCreateTree(t *testing.T) {

}
