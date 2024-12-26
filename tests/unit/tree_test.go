package tests

import (
	"adda/pkg/objects"
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

    dubNode, err := root.GetNode("dub")
    if err != nil {
        t.Error(err.Error())
    }
    
    if dubNode.Path != "dub" {
        t.Errorf("Expected dub; Received: %v", dubNode.Path)
    }
}

func TestAppendNode(t *testing.T) {
    root := objects.NewNode(".")
    root.AppendNode("sub", ".")
    root.AppendNode("dub", ".")
    root.AppendNode("gus", "sub")
    root.AppendNode("foo", "dub")
    root.AppendNode("bar", "gus")
    
    // Test the root's children
    if len(root.Children) == 2 && (root.Children[0].Path != "sub" || root.Children[1].Path != "dub") {
        t.Errorf("Expetced children for root: %v; Received children for root: %v", []string{"sub", "dub"}, []string{root.Children[0].Path, root.Children[1].Path})
    } else if len(root.Children) != 2 {
        t.Errorf("Expected children length: %v; Received children length: %v", 2, len(root.Children))
    }

    // Test the sub's children
    sub, err := root.GetNode("sub")
    if err != nil {
        t.Error(err.Error())
    }

    if len(sub.Children) == 1 && sub.Children[0].Path != "gus" {
        t.Errorf("Expetced children for sub: %v; Received children for sub: %v", []string{"gus"}, []string{sub.Children[0].Path})
    } else if len(sub.Children) != 1 {
        t.Errorf("Expected children length: %v; Received children length: %v", 1, len(sub.Children))
    }

    // Test the dub's children
    dub, err := root.GetNode("dub")
    if err != nil {
        t.Error(err.Error())
    }

    if len(dub.Children) == 1 && dub.Children[0].Path != "foo" {
        t.Errorf("Expetced children for dub: %v; Received children for dub: %v", []string{"foo"}, []string{dub.Children[0].Path})
    } else if len(dub.Children) != 1 {
        t.Errorf("Expected children length: %v; Received children length: %v", 1, len(dub.Children))
    }

    // Test the gus's children
    gus, err := root.GetNode("gus")
    if err != nil {
        t.Error(err.Error())
    }

    if len(gus.Children) == 1 && gus.Children[0].Path != "bar" {
        t.Errorf("Expetced children for gus: %v; Received children for gus: %v", []string{"bar"}, []string{gus.Children[0].Path})
    } else if len(gus.Children) != 1 {
        t.Errorf("Expected children length: %v; Received children length: %v", 1, len(gus.Children))
    }
    
    // Test for foo's children
    foo, err := root.GetNode("foo")
    if err != nil {
        t.Error(err.Error())
    }

    if len(foo.Children) != 0 {
        t.Errorf("Expetced children length: %v; Received children length: %v", 0, len(foo.Children))
    }

    // Test for bar's children
    bar, err := root.GetNode("bar")
    if err != nil {
        t.Error(err.Error())
    }

    if len(bar.Children) != 0 {
        t.Errorf("Expetced children length: %v; Received children length: %v", 0, len(bar.Children))
    }
}

func TestAppendFile(t *testing.T) {

}

func TestContainsChild(t *testing.T) {

}

func TestIndexToTree(t *testing.T) {

}
