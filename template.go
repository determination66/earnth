package earnth

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
)

var _ TemplateEngine = &GoTemplateEngine{}

type TemplateEngine interface {
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
	LoadGlob(pattern string) error
	LoadFromFiles(filenames ...string) error
	LoadFromFS(fs fs.FS, patterns ...string) error
	ParseGlob(pattern string) error
	ParseFiles(filenames ...string) error
	ParseFromFS(fs fs.FS, patterns ...string) error
}

type GoTemplateEngine struct {
	T *template.Template
}

func NewGoTemplateEngine() *GoTemplateEngine {
	return &GoTemplateEngine{}
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

func (g *GoTemplateEngine) LoadGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) LoadFromFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}

func (g *GoTemplateEngine) LoadFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}

// ParseGlob pattern is your render directory
func (g *GoTemplateEngine) ParseGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}

func (g *GoTemplateEngine) ParseFiles(filenames ...string) error {
	var err error
	g.T, err = template.ParseFiles(filenames...)
	return err
}

func (g *GoTemplateEngine) ParseFromFS(fs fs.FS, patterns ...string) error {
	var err error
	g.T, err = template.ParseFS(fs, patterns...)
	return err
}
