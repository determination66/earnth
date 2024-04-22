package earnth

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	Render(ctx context.Context, tplName string, data any) ([]byte, error)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

func (g *GoTemplateEngine) ParseGlob(pattern string) error {
	var err error
	g.T, err = template.ParseGlob(pattern)
	return err
}
