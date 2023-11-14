package mdtemplate

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Parser struct{}

var _ parser.InlineParser = (*Parser)(nil)

var (
	_open  = []byte("{{")
	_close = []byte("}}")
)

// Trigger returns characters that trigger this parser.
func (p *Parser) Trigger() []byte {
	return []byte{'{'}
}
func (p *Parser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	stop := bytes.Index(line, _close)
	if stop < 0 {
		return nil // must close on the same line
	}

	if !bytes.HasPrefix(line, _open) {
		return nil
	}
	segLen := stop + len(_close)
	seg = text.NewSegment(seg.Start, seg.Start+segLen)

	block.Advance(segLen)

	return NewNode(seg)
}
