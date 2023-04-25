package serve

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/abibby/fileserver"
	"github.com/abibby/mdwiki/build"
	"github.com/abibby/mdwiki/util"
)

func Serve(root string, port int) error {
	b, err := build.New(root)
	if err != nil {
		return err
	}
	err = b.Build()
	if err != nil {
		return err
	}

	s := http.NewServeMux()

	s.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		content := r.Form.Get("content")
		file := path.Join(root, util.PathWithoutExt(r.Form.Get("file"))+"md")

		os.WriteFile("")
	})
	s.Handle("/", fileserver.WithFallback(os.DirFS(root), "dist", "index.html", nil))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}
