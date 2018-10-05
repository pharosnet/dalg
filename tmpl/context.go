package tmpl

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
)

var _contextTpl = _notes + `
package %s

import (
	"context"
	"database/sql"
)

const (
	ctxKeyPreparer = "preparer"
)

type Preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func NewContext(parent context.Context, p Preparer) context.Context {
	return context.WithValue(parent, ctxKeyPreparer, p)
}

func prepare(ctx context.Context) Preparer {
	v := ctx.Value(ctxKeyPreparer)
	if v == nil {
		return nil
	}
	return v.(Preparer)
}

`

func WriteContextFile(dbDef *def.Db, dir string) error {
	code := fmt.Sprintf(_contextTpl, dbDef.Package)
	return ioutil.WriteFile(filepath.Join(dir, "context.go"), []byte(code), 0666)
}