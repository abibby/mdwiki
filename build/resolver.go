package build

import (
	"errors"
	"os"
	"path"
	"strings"

	"go.abhg.dev/goldmark/wikilink"
)

type Resolver struct {
	rootDir string
	pages   map[string]bool
}

func NewResolver(dir string) *Resolver {
	return &Resolver{
		rootDir: dir,
		pages:   map[string]bool{},
	}
}

func (r *Resolver) PageExists(name string) (bool, error) {
	exists, ok := r.pages[name]
	if ok {
		return exists, nil
	}
	_, err := os.Stat(path.Join(r.rootDir, name+".md"))
	if errors.Is(err, os.ErrNotExist) {
		r.pages[name] = false
		return false, nil
	} else if err != nil {
		return false, err
	}

	r.pages[name] = true

	return true, nil
}
func (r *Resolver) ResolveWikilink(node *wikilink.Node) ([]byte, error) {
	target := strings.ReplaceAll(strings.ToLower(string(node.Target)), " ", "_")
	exists, err := r.PageExists(string(target))
	if err != nil {
		return nil, err
	}
	fragment := ""
	if !exists {
		fragment = "#create"
	}

	return []byte(target + ".html" + fragment), nil
}
