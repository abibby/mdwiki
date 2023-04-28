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
		"dist/main.js":  res.JS,
	}
	for assetPath, data := range files {
		err := os.WriteFile(path.Join(b.root, assetPath), data, 0644)
		if err != nil {
			return fmt.Errorf("could not copy %s: %w", assetPath, err)
		}
	}
	return nil
}
