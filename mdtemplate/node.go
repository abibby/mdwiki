package mdtemplate

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var Kind = ast.NewNodeKind("mdtemplate")

type Node struct {
	ast.BaseInline

	seg text.Segment
}

var _ ast.Node = (*Node)(nil)

func NewNode(seg text.Segment) *Node {
	return &Node{
		seg: seg,
	}
}

// Kind reports the kind of this node.
func (n *Node) Kind() ast.NodeKind {
	return Kind
}

// Dump dumps the Node to stdout.
func (n *Node) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{}, nil)
}
