package serve

import (
	"fmt"
	"net/http"
	"os"
	"path"

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

	s.Handle("/save", save(root))

	s.Handle("/", http.FileServer(http.Dir(path.Join(root, "dist"))))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}

func save(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		content := r.Form.Get("content")
		file := path.Join(root, util.PathWithoutExt(r.Form.Get("file"))+".md")

		err := os.WriteFile(file, []byte(content), 0644)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed to save page: %v", err)
			return
		}

		b, err := build.New(root)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed to initialize builder: %v", err)
		}

		err = b.BuildFiles([]string{file})
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed to build page: %v", err)
			return
		}

		http.Redirect(w, r, "/"+r.Form.Get("file"), http.StatusFound)
	}
}
