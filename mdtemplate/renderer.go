package mdtemplate

import (
	"fmt"
	"html/template"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type Renderer struct {
	Template *template.Template
	Data     any
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
}

func (r *Renderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*Node)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *mdtemplate.Node", node)
	}

	content := n.seg.Value(src)
	t, err := r.Template.Clone()
	if err != nil {
		return ast.WalkStop, fmt.Errorf("template clone: %w", err)
	}

	t, err = t.Parse(string(content))
	if err != nil {
		return ast.WalkStop, fmt.Errorf("template parse: %w", err)
	}

	err = t.Execute(w, r.Data)
	if err != nil {
		return ast.WalkStop, fmt.Errorf("template execute: %w", err)
	}

	return ast.WalkContinue, nil
}
