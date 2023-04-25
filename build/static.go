package build

import (
	"fmt"
	"os"
	"path"

	"github.com/abibby/mdwiki/res"
)

func (b *Builder) copyStaticFiles() error {
	files := map[string][]byte{
		"dist/main.css": res.CSS,
	}
	for assetPath, data := range files {
		err := os.WriteFile(path.Join(b.root, assetPath), data, 0644)
		if err != nil {
			return fmt.Errorf("could not copy css %w", err)
		}
	}
	return nil
}
