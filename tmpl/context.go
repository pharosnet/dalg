package tmpl

import (
	"bytes"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

var _contextTpl = `
package {{.Package}}

import (
	"context"
	"database/sql"
	_ "{{.Driver}}"
)

const (
	ctxKeyPreparer = "preparer"
	ctxKeyOperator = "operator"
)

type Preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func WithOperator(parent context.Context, op string) context.Context {
	return context.WithValue(parent, ctxKeyOperator, op)
}

func WithPreparer(parent context.Context, p Preparer) context.Context {
	return context.WithValue(parent, ctxKeyPreparer, p)
}

func prepare(ctx context.Context) Preparer {
	v := ctx.Value(ctxKeyPreparer)
	if v == nil {
		return nil
	}
	return v.(Preparer)
}

func operator(ctx context.Context) string {
	v := ctx.Value(ctxKeyOperator)
	if v == nil {
		return ""
	}
	return v.(string)
}

`

func WriteContextFile(dbDef *def.Db, dir string) error {
	tpl, tplErr := template.New("_context").Parse(_contextTpl)
	if tplErr != nil {
		return tplErr
	}
	buffer := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buffer, dbDef); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, "context.go"), buffer.Bytes(), 0666)
}