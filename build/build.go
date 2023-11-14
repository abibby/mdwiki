package build

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/abibby/mdwiki/mdtemplate"
	"github.com/abibby/mdwiki/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/wikilink"
)

type Builder struct {
	root     string
	md       goldmark.Markdown
	template *template.Template
	resolver *Resolver
}

func New(root string) (*Builder, error) {
	t, err := initTemplate(root)
	if err != nil {
		return nil, err
	}
	extenderTemplate, err := t.Clone()
	if err != nil {
		return nil, err
	}

	r := NewResolver(root)
	return &Builder{
		root: root,
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				&wikilink.Extender{
					Resolver: r,
				},
				&mdtemplate.Extender{
					Template: extenderTemplate,
				},
			),
		),
		template: t,
		resolver: r,
	}, nil
}

func (b *Builder) Build() error {
	err := os.RemoveAll(path.Join(b.root, "dist"))
	if err != nil {
		return fmt.Errorf("could not clear dist: %w", err)
	}

	files, err := os.ReadDir(b.root)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
	}

	err = b.mkdirs("dist/edit")
	if err != nil {
		return err
	}

	err = b.copyStaticFiles()
	if err != nil {
		return fmt.Errorf("could not copy static files %w", err)
	}

	inFiles := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() || path.Ext(f.Name()) != ".md" {
			continue
		}
		inFiles = append(inFiles, path.Join(b.root, f.Name()))
	}

	err = b.BuildFiles(inFiles)
	if err != nil {
		return fmt.Errorf("could not build files: %w", err)
	}

	// err = b.updateLinks()
	// if err != nil {
	// 	return fmt.Errorf("update link cache: %w", err)
	// }

	return nil
}

func (b *Builder) mkdirs(paths ...string) error {
	for _, p := range paths {
		err := os.MkdirAll(path.Join(b.root, p), 0755)
		if err != nil {
			return fmt.Errorf("could not create folder %s: %w", p, err)
		}
	}
	return nil
}

func (b *Builder) BuildFiles(inFiles []string) error {
	var err error
	for _, name := range inFiles {
		_, srcFile := path.Split(name)
		err = b.buildFile(
			srcFile,
			path.Join(b.root, "dist", util.PathWithoutExt(srcFile)+".html"),
			path.Join(b.root, "dist/edit", util.PathWithoutExt(srcFile)+".html"),
		)
		if err != nil {
			return fmt.Errorf("counld not build file %s: %w", name, err)
		}
	}

	editPages := []string{}
	for page, exists := range b.resolver.pages {
		if exists {
			continue
		}
		editPages = append(editPages, page+".md")
	}

	for _, srcFile := range editPages {
		err = b.buildFile(
			srcFile,
			path.Join(b.root, "dist", util.PathWithoutExt(srcFile)+".html"),
			path.Join(b.root, "dist/edit", util.PathWithoutExt(srcFile)+".html"),
		)
		if err != nil {
			return fmt.Errorf("counld not build file %s: %w", srcFile, err)
		}
	}

	return nil
}
func (b *Builder) buildFile(in, out, edit string) error {
	b.resolver.SetCurrentPage(in)

	onlyEdit := false
	rawMD, err := os.ReadFile(path.Join(b.root, in))
	if os.IsNotExist(err) {
		rawMD = []byte{}
		onlyEdit = true
	} else if err != nil {
		return fmt.Errorf("read file %s: %w", in, err)
	}

	// t, err := b.template.Clone()
	// if err != nil {
	// 	return fmt.Errorf("clone template: %w", err)
	// }

	// t, err = t.New("default").Parse(string(f))
	// if err != nil {
	// 	return fmt.Errorf("parse template: %w", err)
	// }

	if !onlyEdit {
		// 	srcBuff := &bytes.Buffer{}
		destBuff := &bytes.Buffer{}

		// 	err := t.ExecuteTemplate(srcBuff, in, nil)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to execute src template: %w", err)
		// 	}

		err = b.md.Convert(rawMD, destBuff)
		if err != nil {
			return fmt.Errorf("convert markdown: %w", err)
		}

		// fmt.Println(destBuff.String())
		// t, err = t.New("default").Parse(destBuff.String())
		// if err != nil {
		// 	return fmt.Errorf("parse md template: %w", err)
		// }

		err = b.writeTemplateToFile(out, "default", &DefaultTemplateData{
			Title:       in,
			Body:        template.HTML(destBuff.String()),
			ContentPage: true,
			File:        util.PathWithoutExt(in) + ".html",
		})
		if err != nil {
			return fmt.Errorf("default: %w", err)
		}
	}
	err = b.writeTemplateToFile(edit, "edit", &EditTemplateData{
		Title:   in,
		Content: string(rawMD),
		File:    util.PathWithoutExt(in) + ".html",
	})
	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}
	return nil
}

func (b *Builder) writeTemplateToFile(outFile, templateName string, data any) error {
	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open html file to write: %w", err)
	}
	defer f.Close()
	err = b.template.ExecuteTemplate(f, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
