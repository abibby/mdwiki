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
	s.Handle("/delete", delete(root))

	s.Handle("/", http.FileServer(http.Dir(path.Join(root, "dist"))))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), s)
}

func save(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		content := r.Form.Get("content")
		file := path.Join(root, util.PathWithoutExt(r.Form.Get("file"))+".md")

		err := os.WriteFile(file, []byte(content), 0644)
		if checkError(w, err) {
			return
		}

		err = updatePages(root)
		if checkError(w, err) {
			return
		}

		http.Redirect(w, r, "/"+r.Form.Get("file"), http.StatusFound)
	}
}

func delete(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file := path.Join(root, util.PathWithoutExt(r.Form.Get("file"))+".md")

		err := os.Remove(file)
		if checkError(w, err) {
			return
		}

		err = updatePages(root)
		if checkError(w, err) {
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func updatePages(root string) error {
	b, err := build.New(root)
	if err != nil {
		return fmt.Errorf("failed to initialize builder: %w", err)
	}

	err = b.Build()
	// err = b.BuildFiles([]string{file})
	if err != nil {
		return fmt.Errorf("failed to build page: %w", err)
	}
	return nil
}

func checkError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))

	return true
}
