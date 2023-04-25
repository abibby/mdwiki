package build

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (b *Builder) updateLinks() error {
	linkRoot := path.Join(b.root, "links")
	links := b.resolver.Links()
	for from, tos := range links {
		b, err := os.ReadFile(path.Join(linkRoot, from))
		if os.IsNotExist(err) {
			b = []byte{}
		} else if err != nil {
			return err
		}
		_ = strings.Split(string(b), "\n")
		for _, to := range tos {
			f, err := os.OpenFile(path.Join(linkRoot, to), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return fmt.Errorf("failed to open link file %s: %w", to, err)
			}
			_, err = fmt.Fprintln(f, from)
			if err != nil {
				return fmt.Errorf("failed to update link file %s: %w", to, err)
			}
		}
	}
	return nil
}
