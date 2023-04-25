package build

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path"
	texttemplate "text/template"

	"github.com/abibby/mdwiki/res"
	"github.com/abibby/mdwiki/util"
)

type DefaultTemplateData struct {
	Title       string
	Body        template.HTML
	ContentPage bool
	File        string
}

type EditTemplateData struct {
	Title   string
	Content string
	File    string
}

func initTemplate(root string) (*template.Template, error) {
	t := template.New("default")

	t, err := addFuncs(t, root)
	if err != nil {
		return nil, fmt.Errorf("could not add functions: %w", err)
	}
	t, err = t.Parse(res.HtmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}
	t, err = addSubTemplates(t, map[string]string{
		"edit": res.Edit,
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}

	return t, nil
}

func addSubTemplates(t *template.Template, templates map[string]string) (*template.Template, error) {
	defaultTemplate, err := texttemplate.New("default").Parse(res.HtmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to clone default template: %w", err)
	}
	defaultTemplate.Templates()
	for name, data := range templates {
		buff := &bytes.Buffer{}
		defaultTemplate.Execute(buff, &DefaultTemplateData{
			Title: "{{ .Title }}",
			File:  "{{ .File }}",
			Body:  template.HTML(data),
		})

		t, err = t.New(name).Parse(buff.String())
		if err != nil {
			return nil, fmt.Errorf("could not parse template: %w", err)
		}
	}

	return t, nil
}

func addFuncs(t *template.Template, dir string) (*template.Template, error) {
	root := path.Join(dir, "functions")
	files, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return t, nil
	} else if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{}

	for _, f := range files {
		b, err := os.ReadFile(path.Join(root, f.Name()))
		if err != nil {
			return nil, err
		}
		funcs[util.PathWithoutExt(f.Name())] = func(args ...any) template.HTML {
			return template.HTML(fmt.Sprintf(string(b), args...))
		}
	}

	return t.Funcs(funcs), nil
}
