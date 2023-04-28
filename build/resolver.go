package build

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/abibby/mdwiki/util"
	"go.abhg.dev/goldmark/wikilink"
)

type Resolver struct {
	rootDir     string
	pages       map[string]bool
	currentPage string
	links       map[string][]string
}

func NewResolver(dir string) *Resolver {
	return &Resolver{
		rootDir: dir,
		pages:   map[string]bool{},
		links:   map[string][]string{},
	}
}

func (r *Resolver) addLink(from, to string) {
	p, ok := r.links[from]
	if !ok {
		p = []string{to}
	} else {
		p = append(p, to)
	}
	r.links[from] = p
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

func (r *Resolver) SetCurrentPage(p string) {
	r.currentPage = p
}

func (r *Resolver) Links() map[string][]string {
	return r.links
}

func (r *Resolver) ResolveWikilink(node *wikilink.Node) ([]byte, error) {
	if node.Embed {
		return node.Target, nil
	}
	target := strings.ReplaceAll(strings.ToLower(string(node.Target)), " ", "_")
	exists, err := r.PageExists(string(target))
	if err != nil {
		return nil, err
	}
	r.addLink(util.PathWithoutExt(r.currentPage), string(target))

	if !exists {
		return []byte(path.Join("/edit", target+".html")), nil
	}

	return []byte("/" + target + ".html"), nil
}
