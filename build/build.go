package build

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/abibby/mdwiki/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/wikilink"
)

type Builder struct {
	root     string
	md       goldmark.Markdown
	template *template.Template
}

func New(root string) (*Builder, error) {
	t, err := initTemplate(root)
	if err != nil {
		return nil, err
	}
	return &Builder{
		root: root,
		md: goldmark.New(
			goldmark.WithExtensions(extension.GFM, &wikilink.Extender{
				Resolver: NewResolver(root),
			}),
		),
		template: t,
	}, nil
}

func (b *Builder) Build() error {
	files, err := os.ReadDir(b.root)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
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

	return b.BuildFiles(inFiles)
}

func (b *Builder) BuildFiles(inFiles []string) error {
	_, err := b.template.ParseFiles(inFiles...)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %w", err)
	}
	err = os.MkdirAll(path.Join(b.root, "dist/edit"), 0755)
	if err != nil {
		return fmt.Errorf("could not create destination folder %w", err)
	}

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

	return nil
}
func (b *Builder) buildFile(in, out, edit string) error {
	srcBuff := &bytes.Buffer{}
	destBuff := &bytes.Buffer{}

	err := b.template.ExecuteTemplate(srcBuff, in, nil)
	if err != nil {
		return fmt.Errorf("failed to execute src template: %w", err)
	}

	err = b.md.Convert(srcBuff.Bytes(), destBuff)
	if err != nil {
		return fmt.Errorf("failed convert markdown: %w", err)
	}

	err = b.writeTemplateToFile(out, "default", &DefaultTemplateData{
		Title:       in,
		Body:        template.HTML(destBuff.String()),
		ContentPage: true,
		File:        util.PathWithoutExt(in) + ".html",
	})
	if err != nil {
		return fmt.Errorf("default: %w", err)
	}
	rawMD, err := os.ReadFile(path.Join(b.root, in))
	if err != nil {
		return fmt.Errorf("failed to read markdown for editing: %w", err)
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

func (b *Builder) writeTemplateToFile(file, templateName string, data any) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
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
